package main

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/builder"
	"github.com/nnurry/harmonia/internal/service"
)

func main() {
	fmt.Println("Opening conn")

	// https://libvirt.org/uri.html#keyfile-parameter
	connBuilder := builder.
		NewLibvirtConnectBuilder().
		WithTransportType("ssh").
		WithUser("nnurry").
		WithHost("xeon-opensuse").
		WithKeyfilePath("/develop/.host-ssh/xeon-opensuse")

	libvirtService, err := service.NewLibvirtFromConnectBuilder(connBuilder)

	if err != nil {
		panic(err)
	}

	fmt.Println("Opened conn")

	// TODO: make these deps useful
	baseDomainName := "leap-base-VM-latest-test"
	newDomainName := "leap-base-VM-latest-test-new"

	qcow2DiskPath := "/var/lib/libvirt/images/leap-base-VM-latest-test-new.qcow2"
	cloudInitDiskPath := "/var/lib/libvirt/images/leap-base-VM-latest-test-new.iso"

	baseDomain, err := libvirtService.GetDomainByName(baseDomainName)
	if err != nil {
		panic(err)
	}

	numOfCpus := 4
	memory := uint(8 * 1024 * 1024)
	newDomainBuilder, err := builder.NewLibvirtDomainBuilder(baseDomain)
	if err != nil {
		panic(err)
	}

	newDomainBuilder = newDomainBuilder.
		WithDomainName(newDomainName).
		WithNumOfCpus(numOfCpus).
		WithMemory(memory, "KiB").
		WithCiDiskPath(cloudInitDiskPath).
		WithQcow2DiskPath(qcow2DiskPath)

	dom, err := libvirtService.DefineDomainFromBuilder(newDomainBuilder)
	if err != nil {
		panic(err)
	}

	err = dom.Create()
	if err != nil {
		panic(err)
	}

	fmt.Printf("VM %v successfully defined\n", newDomainName)

	fmt.Println("Closing conn")
	libvirtService.Cleanup()
	fmt.Println("Closed conn")
}
