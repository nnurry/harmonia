package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nnurry/harmonia/internal/builder"
	"github.com/nnurry/harmonia/internal/connection"
	"github.com/nnurry/harmonia/internal/contract"
	"github.com/nnurry/harmonia/internal/logger"
	"github.com/nnurry/harmonia/internal/processor"
	"github.com/nnurry/harmonia/internal/service/cloudinit"
	"github.com/nnurry/harmonia/pkg/utils"
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

func NewVirtualMachineFromVirtualMachineConfig(config contract.VirtualMachineConfig) (*VirtualMachine, error) {
	var (
		sshConnection    *connection.SSH
		shellProcessor   ShellProcessor
		libvirtService   LibvirtService
		cloudInitService CloudInitService
	)

	if config.HypervisorConnectionConfig.IsLocalShell {
		shellProcessor = processor.NewLocalShell()
	} else {
		var err error
		sshConnection, err = connection.NewSSH(config.HypervisorConnectionConfig.SSHConfig)
		if err != nil {
			return nil, err
		}
		shellProcessor = processor.NewSecureShell(sshConnection)
	}

	// create services
	if conn, err := connection.NewLibvirt(config.HypervisorConnectionConfig.LibvirtConfig); err != nil {
		return nil, err
	} else {
		libvirtService, err = NewLibvirt(conn)
		if err != nil {
			return nil, err
		}
	}

	cloudInitService, err := NewCloudInit(shellProcessor, sshConnection)
	if err != nil {
		return nil, err
	}

	return NewVirtualMachine(libvirtService, cloudInitService, shellProcessor)
}

func (service *VirtualMachine) Create(config contract.VirtualMachineConfig) (string, error) {
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
					MacAddress:         config.NetworkVMConfig.MacAddress,
					Nameservers:        cloudinit.Nameserver{Addresses: config.NetworkVMConfig.Nameservers},
				},
			},
		},
	})

	cloudInitDir := fmt.Sprintf("/var/my-cloud-init/%v/%v", config.GeneralVMConfig.Name, uniqueID)

	logger.Info("creating cloud-init.iso")
	cloudInitIsoPath, err := service.cloudInitService.WriteToDisk(context.Background(), cloudInitDir, "cloud-init.iso")
	if err != nil {
		return "", err
	}
	logger.Info("created cloud-init.iso")
	defer service.Cleanup(cloudInitDir)

	// create VM
	logger.Infof("creating libvirt domain from %v\n", config.GeneralVMConfig.BaseVirtualMachineName)

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

	baseQCOW2Path := baseQCOW2Disk.Source.File.File
	baseQCOW2PathAsParts := strings.Split(baseQCOW2Path, "/")

	newQCOW2Path := fmt.Sprintf(
		"%v/%v.qcow2",
		strings.Join(baseQCOW2PathAsParts[:len(baseQCOW2PathAsParts)-1], "/"),
		config.GeneralVMConfig.Name,
	)
	if err = service.cloneDisk(baseQCOW2Path, newQCOW2Path, config.DiskSizeInGiB, config.IsCopyOnWriteClone); err != nil {
		service.revertCloudInitChange <- true
		return "", err
	}

	libvirtBuilder, err := builder.NewLibvirtDomainBuilder(
		baseDomain,
		[]*builder.DomainBuilderFlag{builder.SET_VM_NAME}, // no need any flag
		false,
	)
	if err != nil {
		logger.Infof("failed to create libvirt domain builder: %v\n", err)
		service.revertCloudInitChange <- true
		return "", err
	}

	libvirtBuilder = libvirtBuilder.
		WithDomainName(config.GeneralVMConfig.Name).
		WithCiDiskPath(cloudInitIsoPath).
		WithQcow2DiskPath(newQCOW2Path).
		WithMemory(uint(config.GeneralVMConfig.MemoryInGiB*1024*1024), "KiB").
		WithNumOfCpus(config.GeneralVMConfig.NumOfVCPUs).
		WithMacAddress(config.NetworkVMConfig.MacAddress)

	newDomain, err := service.libvirtService.DefineDomainFromBuilder(libvirtBuilder)
	if err != nil {
		service.revertCloudInitChange <- true
		return "", err
	}
	logger.Info("created libvirt domain")
	logger.Info("starting VM")
	if err = newDomain.Create(); err != nil {
		return "", fmt.Errorf("failed to start VM")
	}
	logger.Info("started VM")

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
		logger.Info("removed cloud-init iso after failing to create VM")
	}
	return nil
}

func (service *VirtualMachine) Delete(config contract.VirtualMachineConfig) (string, error) {
	return "", fmt.Errorf("unimplemented")
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
	stderrBuffer := bytes.NewBuffer([]byte{})

	err := service.shellProcessor.Execute(
		context.Background(),
		os.Stdout, stderrBuffer,
		command, arguments...,
	)

	if err != nil {
		stdErrAsString := stderrBuffer.String()
		logger.Error(stdErrAsString)

		err = fmt.Errorf("could not clone disk: %v", stdErrAsString)
		return err
	}
	return nil
}
