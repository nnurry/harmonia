package service

import (
	"fmt"

	"libvirt.org/go/libvirtxml"
)

type VmBuilderFlag string

const (
	SET_VM_NAME         = VmBuilderFlag("set VM name")
	SET_NUM_OF_CPUS     = VmBuilderFlag("set num of CPUs")
	SET_MEMORY          = VmBuilderFlag("set memory")
	SET_QCOW2_DISK_PATH = VmBuilderFlag("set qcow2 disk path")
	SET_CI_DISK_PATH    = VmBuilderFlag("set cloud-init disk path")
)

type LibvirtDomainBuilder struct {
	baseVmXml       *libvirtxml.Domain
	newVmXml        *libvirtxml.Domain
	qcow2DomainDisk *libvirtxml.DomainDisk
	ciDomainDisk    *libvirtxml.DomainDisk

	requiredFlags map[VmBuilderFlag]bool
}

func NewLibvirtDomainBuilder(baseVmXml *libvirtxml.Domain) (*LibvirtDomainBuilder, error) {
	if baseVmXml == nil {
		return nil, fmt.Errorf("can't create vm builder: base vm xml object is nil")
	}

	builder := &LibvirtDomainBuilder{
		baseVmXml: baseVmXml,
		requiredFlags: map[VmBuilderFlag]bool{
			SET_VM_NAME:         false,
			SET_NUM_OF_CPUS:     false,
			SET_MEMORY:          false,
			SET_CI_DISK_PATH:    false,
			SET_QCOW2_DISK_PATH: false,
		},
	}

	newVmXml := &libvirtxml.Domain{}
	newVmXml.Type = "kvm"
	newVmXml.Metadata = baseVmXml.Metadata
	newVmXml.OS = baseVmXml.OS

	newVmXml.Features = baseVmXml.Features
	newVmXml.Clock = baseVmXml.Clock
	newVmXml.OnPoweroff = baseVmXml.OnPoweroff
	newVmXml.OnReboot = baseVmXml.OnReboot
	newVmXml.OnCrash = baseVmXml.OnCrash

	newVmXml.PM = &libvirtxml.DomainPM{
		SuspendToMem:  &libvirtxml.DomainPMPolicy{Enabled: "no"},
		SuspendToDisk: &libvirtxml.DomainPMPolicy{Enabled: "no"},
	}

	newVmXml.Devices = baseVmXml.Devices

	builder.newVmXml = newVmXml
	return builder, nil
}
func (builder *LibvirtDomainBuilder) WithVmName(name string) *LibvirtDomainBuilder {
	builder.newVmXml.Name = name

	builder.requiredFlags[SET_VM_NAME] = true
	return builder
}

func (builder *LibvirtDomainBuilder) WithNumOfCpus(numOfCpus int) *LibvirtDomainBuilder {
	builder.newVmXml.CPU = &libvirtxml.DomainCPU{
		Mode: builder.baseVmXml.CPU.Mode,
		Topology: &libvirtxml.DomainCPUTopology{
			Sockets: numOfCpus,
			Threads: 1,
			Cores:   1,
		},
	}
	builder.newVmXml.VCPU = &libvirtxml.DomainVCPU{
		Placement: "static",
		Current:   uint(numOfCpus),
		Value:     uint(numOfCpus),
	}

	builder.requiredFlags[SET_NUM_OF_CPUS] = true
	return builder
}

func (builder *LibvirtDomainBuilder) WithMemory(memory uint, unit string) *LibvirtDomainBuilder {
	builder.newVmXml.Memory = &libvirtxml.DomainMemory{Value: memory, Unit: unit}
	builder.newVmXml.CurrentMemory = &libvirtxml.DomainCurrentMemory{Value: memory, Unit: unit}

	builder.requiredFlags[SET_MEMORY] = true
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

	builder.requiredFlags[SET_QCOW2_DISK_PATH] = true
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

	builder.requiredFlags[SET_CI_DISK_PATH] = true
	return builder
}

func (builder *LibvirtDomainBuilder) BuildXMLString() (string, error) {
	unsatisfiedFlags := []VmBuilderFlag{}
	for flag, satisfied := range builder.requiredFlags {
		if !satisfied {
			unsatisfiedFlags = append(unsatisfiedFlags, flag)
		}
	}

	if len(unsatisfiedFlags) > 0 {
		return "", fmt.Errorf("failed to build VM's XML: flag [%v] not satisfied", unsatisfiedFlags)
	}

	builder.newVmXml.Devices.Disks = []libvirtxml.DomainDisk{
		*builder.qcow2DomainDisk,
		*builder.ciDomainDisk,
	}

	xmlString, err := builder.newVmXml.Marshal()
	if err != nil {
		return "", fmt.Errorf("failed to build VM's xml: can't serialize to XML string due to %v", err)
	}

	return xmlString, nil
}
