package builder

import (
	"fmt"

	"github.com/nnurry/harmonia/pkg/types"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type DomainBuilderFlag types.BuilderFlag

const (
	SET_VM_NAME         = DomainBuilderFlag("set VM name")
	SET_NUM_OF_CPUS     = DomainBuilderFlag("set num of CPUs")
	SET_MEMORY          = DomainBuilderFlag("set memory")
	SET_QCOW2_DISK_PATH = DomainBuilderFlag("set qcow2 disk path")
	SET_CI_DISK_PATH    = DomainBuilderFlag("set cloud-init disk path")
)

type LibvirtDomainBuilder struct {
	baseDomainXml   *libvirtxml.Domain
	newDomainXml    *libvirtxml.Domain
	qcow2DomainDisk *libvirtxml.DomainDisk
	ciDomainDisk    *libvirtxml.DomainDisk

	requiredFlagMap map[DomainBuilderFlag]bool
}

func NewLibvirtDomainBuilder(baseDomain *libvirt.Domain, requiredFlags ...DomainBuilderFlag) (*LibvirtDomainBuilder, error) {
	baseDomainXmlString, err := baseDomain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	if err != nil {
		return nil, fmt.Errorf("unable to get base domain's XML string definition: %v", err)
	}

	baseDomainXml := &libvirtxml.Domain{}
	err = baseDomainXml.Unmarshal(baseDomainXmlString)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize base domain's XML string definition: %v", err)
	}

	newDomainXml := &libvirtxml.Domain{}
	newDomainXml.Type = "domain"
	newDomainXml.Metadata = baseDomainXml.Metadata
	newDomainXml.OS = baseDomainXml.OS

	newDomainXml.Features = baseDomainXml.Features
	newDomainXml.Clock = baseDomainXml.Clock
	newDomainXml.OnPoweroff = baseDomainXml.OnPoweroff
	newDomainXml.OnReboot = baseDomainXml.OnReboot
	newDomainXml.OnCrash = baseDomainXml.OnCrash

	newDomainXml.PM = &libvirtxml.DomainPM{
		SuspendToMem:  &libvirtxml.DomainPMPolicy{Enabled: "no"},
		SuspendToDisk: &libvirtxml.DomainPMPolicy{Enabled: "no"},
	}

	newDomainXml.Devices = baseDomainXml.Devices

	if len(requiredFlags) < 1 {
		requiredFlags = []DomainBuilderFlag{
			SET_VM_NAME,
			SET_NUM_OF_CPUS,
			SET_MEMORY,
			SET_CI_DISK_PATH,
			SET_QCOW2_DISK_PATH,
		}
	}

	requiredFlagMap := make(map[DomainBuilderFlag]bool)
	for _, flag := range requiredFlags {
		requiredFlagMap[flag] = false
	}

	builder := &LibvirtDomainBuilder{
		baseDomainXml:   baseDomainXml,
		newDomainXml:    newDomainXml,
		requiredFlagMap: requiredFlagMap,
	}
	return builder, nil
}
func (builder *LibvirtDomainBuilder) WithDomainName(name string) *LibvirtDomainBuilder {
	builder.newDomainXml.Name = name

	builder.requiredFlagMap[SET_VM_NAME] = true
	return builder
}

func (builder *LibvirtDomainBuilder) WithNumOfCpus(numOfCpus int) *LibvirtDomainBuilder {
	builder.newDomainXml.CPU = &libvirtxml.DomainCPU{
		Mode: builder.baseDomainXml.CPU.Mode,
		Topology: &libvirtxml.DomainCPUTopology{
			Sockets: numOfCpus,
			Threads: 1,
			Cores:   1,
		},
	}
	builder.newDomainXml.VCPU = &libvirtxml.DomainVCPU{
		Placement: "static",
		Current:   uint(numOfCpus),
		Value:     uint(numOfCpus),
	}

	builder.requiredFlagMap[SET_NUM_OF_CPUS] = true
	return builder
}

func (builder *LibvirtDomainBuilder) WithMemory(memory uint, unit string) *LibvirtDomainBuilder {
	builder.newDomainXml.Memory = &libvirtxml.DomainMemory{Value: memory, Unit: unit}
	builder.newDomainXml.CurrentMemory = &libvirtxml.DomainCurrentMemory{Value: memory, Unit: unit}

	builder.requiredFlagMap[SET_MEMORY] = true
	return builder
}

func (builder *LibvirtDomainBuilder) WithQcow2DiskPath(path string) *LibvirtDomainBuilder {
	builder.qcow2DomainDisk = &libvirtxml.DomainDisk{
		Device: "disk",
		Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "qcow2", Cache: "none", Discard: "unmap"},
		Target: &libvirtxml.DomainDiskTarget{Dev: "vdb", Bus: "virtio"},
		Source: &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: path}},
		Boot:   &libvirtxml.DomainDeviceBoot{Order: 1},
	}

	builder.requiredFlagMap[SET_QCOW2_DISK_PATH] = true
	return builder
}

func (builder *LibvirtDomainBuilder) WithCiDiskPath(path string) *LibvirtDomainBuilder {
	builder.ciDomainDisk = &libvirtxml.DomainDisk{
		Device:   "cdrom",
		Driver:   &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "raw"},
		Source:   &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: path}},
		Target:   &libvirtxml.DomainDiskTarget{Dev: "hdc", Bus: "sata"},
		ReadOnly: &libvirtxml.DomainDiskReadOnly{},
		Boot:     &libvirtxml.DomainDeviceBoot{Order: 2},
	}

	builder.requiredFlagMap[SET_CI_DISK_PATH] = true
	return builder
}

func (builder *LibvirtDomainBuilder) BuildXMLString() (string, error) {
	unsatisfiedFlags := []DomainBuilderFlag{}
	for flag, satisfied := range builder.requiredFlagMap {
		if !satisfied {
			unsatisfiedFlags = append(unsatisfiedFlags, flag)
		}
	}

	if len(unsatisfiedFlags) > 0 {
		return "", fmt.Errorf("flag [%v] not satisfied", unsatisfiedFlags)
	}

	builder.newDomainXml.Devices.Disks = []libvirtxml.DomainDisk{
		*builder.qcow2DomainDisk,
		*builder.ciDomainDisk,
	}

	xmlString, err := builder.newDomainXml.Marshal()
	if err != nil {
		return "", fmt.Errorf("unable to serialize domain from XML definition due to %v", err)
	}

	return xmlString, nil
}

func (builder *LibvirtDomainBuilder) Build(conn *libvirt.Connect) (*libvirt.Domain, error) {
	xmlDefinition, err := builder.BuildXMLString()
	if err != nil {
		return nil, err
	}

	domain, err := conn.DomainDefineXML(xmlDefinition)
	if err != nil {
		return nil, err
	}

	return domain, err
}
