package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nnurry/harmonia/internal/builder"
	"github.com/nnurry/harmonia/internal/contract"
	"github.com/nnurry/harmonia/internal/service/cloudinit"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/rs/zerolog/log"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type VirtualMachine struct {
	libvirtService        LibvirtService
	cloudInitService      CloudInitService
	shellProcessor        ShellProcessor
	revertCloudInitChange chan bool
}

func NewVirtualMachine(
	libvirtService LibvirtService,
	cloudInitService CloudInitService,
	shellProcessor ShellProcessor) (*VirtualMachine, error) {

	return &VirtualMachine{
		libvirtService:        libvirtService,
		cloudInitService:      cloudInitService,
		shellProcessor:        shellProcessor,
		revertCloudInitChange: make(chan bool, 1),
	}, nil

}

func (service *VirtualMachine) Create(config contract.BuildVirtualMachineConfig) (string, error) {
	uniqueID := utils.GenerateUniqueTimestamp()

	// make sure base VM already exists
	baseDomain, err := service.libvirtService.GetDomainByName(config.GeneralVMConfig.BaseVirtualMachineName)
	if err != nil {
		return "", err
	}

	// create cloud-init ISO
	service.cloudInitService.SetMetaData(cloudinit.MetaData{
		Hostname:   config.GeneralVMConfig.Name,
		InstanceId: fmt.Sprintf("%v-%v", config.GeneralVMConfig.Name, uniqueID),
	})
	service.cloudInitService.SetUserData(cloudinit.UserData{
		Hostname:       config.GeneralVMConfig.Name,
		ManageEtcHosts: true,
		DisableRootPw:  true,
		Users: []cloudinit.User{{
			Name:           config.UserVMConfig.User,
			Sudo:           "ALL=(ALL) NOPASSWD:ALL",
			AuthorizedKeys: config.AuthorizedKeyContents,
		}},
	})
	service.cloudInitService.SetNetworkConfig(cloudinit.NetworkConfig{
		Network: cloudinit.Network{
			Version: 2,
			Ethernets: cloudinit.Ethernet{
				Eth0: cloudinit.Eth0{
					Dhcp4:              false,
					IPv4Addresses:      []string{fmt.Sprintf("%v/24", config.NetworkVMConfig.IPv4Address)},
					IPv4GatewayAddress: config.NetworkVMConfig.IPv4GatewayAddress,
					Nameservers:        []cloudinit.Nameserver{{Addresses: config.NetworkVMConfig.Nameservers}},
				},
			},
		},
	})

	cloudInitDir := fmt.Sprintf("/tmp/%v/%v", config.GeneralVMConfig.Name, uniqueID)

	log.Info().Msg("creating cloud-init.iso")
	cloudInitIsoPath, err := service.cloudInitService.WriteToDisk(context.Background(), cloudInitDir, "cloud-init.iso")
	if err != nil {
		return "", err
	}
	log.Info().Msg("created cloud-init.iso")
	defer service.Cleanup(cloudInitDir)

	// create VM
	log.Info().Msgf("creating libvirt domain from %v\n", config.GeneralVMConfig.BaseVirtualMachineName)

	baseDomainXMLDesc, err := baseDomain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	if err != nil {
		service.revertCloudInitChange <- true
		return "", err
	}
	baseDomainXML := &libvirtxml.Domain{}
	err = baseDomainXML.Unmarshal(baseDomainXMLDesc)
	if err != nil {
		service.revertCloudInitChange <- true
		return "", err
	}

	var baseQCOW2Disk *libvirtxml.DomainDisk
	for _, disk := range baseDomainXML.Devices.Disks {
		if disk.Device == "disk" && disk.Driver.Type == "qcow2" {
			baseQCOW2Disk = &disk
			break
		}
	}

	if baseQCOW2Disk == nil {
		service.revertCloudInitChange <- true
		return "", fmt.Errorf("could not get QCOW2 disk from base VM %v", config.BaseVirtualMachineName)
	}

	basePath := baseQCOW2Disk.Source.File.File
	basePathAsParts := strings.Split(basePath, "/")

	newPath := fmt.Sprintf(
		"%v/%v.qcow2",
		strings.Join(basePathAsParts[:len(basePathAsParts)-1], "/"),
		config.GeneralVMConfig.Name,
	)
	if err = service.cloneDisk(basePath, newPath, config.DiskSizeInGiB, config.IsCopyOnWriteClone); err != nil {
		service.revertCloudInitChange <- true
		return "", err
	}

	libvirtBuilder, err := builder.NewLibvirtDomainBuilder(
		baseDomain,
		[]*builder.DomainBuilderFlag{builder.SET_VM_NAME}, // no need any flag
		false,
	)
	if err != nil {
		log.Info().Msgf("failed to create libvirt domain builder: %v\n", err)
		service.revertCloudInitChange <- true
		return "", err
	}

	log.Info().Msgf("config: %v", config)
	libvirtBuilder = libvirtBuilder.
		WithDomainName(config.GeneralVMConfig.Name).
		WithCiDiskPath(cloudInitIsoPath).
		WithMemory(uint(config.GeneralVMConfig.MemoryInGiB*1024*1024), "KiB").
		WithNumOfCpus(config.GeneralVMConfig.NumOfVCPUs)

	newDomain, err := service.libvirtService.DefineDomainFromBuilder(libvirtBuilder)
	if err != nil {
		service.revertCloudInitChange <- true
		return "", err
	}
	log.Info().Msg("created libvirt domain")
	log.Info().Msg("starting VM")
	if err = newDomain.Create(); err != nil {
		return "", fmt.Errorf("failed to start VM")
	}
	log.Info().Msg("started VM")

	service.revertCloudInitChange <- false
	return newDomain.GetUUIDString()
}

func (service *VirtualMachine) Cleanup(cloudInitDir string) error {
	// revert cloud-init change
	if <-service.revertCloudInitChange {
		err := service.cloudInitService.RemoveFromDisk(context.Background(), cloudInitDir)
		if err != nil {
			return fmt.Errorf("failed to remove cloud-init iso after failing to create VM: %v", err)
		}
	}
	return nil
}

func (service *VirtualMachine) cloneDisk(basePath, newPath string, diskSizeInGiB float64, isCopyOnWrite bool) error {
	var (
		command   string
		arguments []string
	)
	if isCopyOnWrite {
		command = "qemu-img"
		arguments = []string{
			"create",
			"-f", "qcow2",
			"-b", basePath,
			"-F", "qcow2",
			newPath,
			fmt.Sprintf("%02fG", diskSizeInGiB),
		}
	} else {
		command = "cp"
		arguments = []string{
			basePath,
			newPath,
		}
	}
	err := service.shellProcessor.Execute(
		context.Background(),
		os.Stdout, os.Stderr,
		command, arguments...,
	)
	if err != nil {
		return err
	}
	return nil
}
