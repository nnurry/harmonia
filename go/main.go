package main

import (
	"fmt"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func main() {
	fmt.Println("Opening conn")
	conn, err := libvirt.NewConnect("qemu+ssh://xeon-opensuse/system")
	if err != nil {
		panic(err)
	}
	fmt.Println("Opened conn")
	baseVmName := "leap-base-VM-latest-test"
	newVmName := "leap-base-VM-latest-test-new"

	qcow2DiskPath := "/var/lib/libvirt/images/leap-base-VM-latest-test-new.qcow2"
	cloudInitDiskPath := "/var/lib/libvirt/images/leap-base-VM-latest-test-new.iso"

	baseVM, err := conn.LookupDomainByName(baseVmName)
	if err != nil {
		panic(err)
	}

	baseVmXmlStr, err := baseVM.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	if err != nil {
		panic(err)
	}

	baseVmXml := &libvirtxml.Domain{}
	err = baseVmXml.Unmarshal(baseVmXmlStr)
	if err != nil {
		panic(err)
	}

	numOfCpus := 4
	memory := uint(4 * 1024 * 1024)
	newVmXml := libvirtxml.Domain{}

	newVmXml.Name = newVmName
	newVmXml.Metadata = baseVmXml.Metadata
	newVmXml.OS = baseVmXml.OS

	newVmXml.Memory = &libvirtxml.DomainMemory{Value: memory, Unit: "KiB"}
	newVmXml.CurrentMemory = &libvirtxml.DomainCurrentMemory{Value: memory, Unit: "KiB"}
	newVmXml.CPU = &libvirtxml.DomainCPU{
		Mode: baseVmXml.CPU.Mode,
		Topology: &libvirtxml.DomainCPUTopology{
			Sockets: numOfCpus,
			Threads: 1,
			Cores:   1,
		},
	}
	newVmXml.VCPU = &libvirtxml.DomainVCPU{
		Placement: "static",
		Current:   uint(numOfCpus),
		Value:     uint(numOfCpus),
	}

	newVmXml.Features = baseVmXml.Features
	newVmXml.Clock = baseVmXml.Clock
	newVmXml.OnPoweroff = baseVmXml.OnPoweroff
	newVmXml.OnReboot = baseVmXml.OnReboot
	newVmXml.OnCrash = baseVmXml.OnCrash
	newVmXml.PM = &libvirtxml.DomainPM{
		SuspendToMem:  &libvirtxml.DomainPMPolicy{Enabled: "no"},
		SuspendToDisk: &libvirtxml.DomainPMPolicy{Enabled: "no"},
	}

	// 1) clone base VM disk
	// 2) add cloud-init disk
	qcow2Disk := libvirtxml.DomainDisk{
		Device: "disk",
		Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "qcow2", Cache: "none", Discard: "unmap"},
		Target: &libvirtxml.DomainDiskTarget{Dev: "vdb", Bus: "virtio"},
		Source: &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: qcow2DiskPath}},
		Boot:   &libvirtxml.DomainDeviceBoot{Order: 1},
	}
	cloudInitDisk := libvirtxml.DomainDisk{
		Device:   "cdrom",
		Driver:   &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "raw"},
		Source:   &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: cloudInitDiskPath}},
		Target:   &libvirtxml.DomainDiskTarget{Dev: "hdc", Bus: "sata"},
		ReadOnly: &libvirtxml.DomainDiskReadOnly{},
		Boot:     &libvirtxml.DomainDeviceBoot{Order: 2},
	}
	newVmXml.Devices = baseVmXml.Devices
	newVmXml.Devices.Disks = []libvirtxml.DomainDisk{qcow2Disk, cloudInitDisk}

	newVmXmlStr, err := baseVmXml.Marshal()
	if err != nil {
		panic(err)
	}
	dom, err := conn.DomainDefineXML(newVmXmlStr)
	if err != nil {
		panic(err)
	}

	err = dom.Create()
	if err != nil {
		panic(err)
	}

	fmt.Printf("VM %v successfully defined\n", newVmName)

	fmt.Println("Closing conn")
	conn.Close()
	fmt.Println("Closed conn")
}
