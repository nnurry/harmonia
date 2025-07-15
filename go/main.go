package main

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/service"
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
	memory := uint(8 * 1024 * 1024)
	newVmBuilder, err := service.NewLibvirtVmBuilder(baseVmXml)
	if err != nil {
		panic(err)
	}

	newVmXmlStr, err := newVmBuilder.
		SetVmName(newVmName).
		SetNumOfCpus(numOfCpus).
		SetMemory(memory, "KiB").
		SetCiDiskPath(cloudInitDiskPath).
		SetQcow2DiskPath(qcow2DiskPath).
		BuildXMLString()

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
