package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nnurry/harmonia/cmd"
	"github.com/nnurry/harmonia/internal/homevirt/builder"
	"github.com/nnurry/harmonia/internal/homevirt/service"
	"github.com/urfave/cli/v2"
)

func test() {
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

func main() {
	cliCommands := &cli.Command{
		Name:        "cli",
		Description: "Commands for interacting with Harmonia's features directly.",
		Subcommands: []*cli.Command{
			{
				Name:        "libvirt",
				Description: "Interact with Libvirt via Go CLI",
				Subcommands: []*cli.Command{
					cmd.NewDefineLibvirtDomainCommand(),
				},
			},
		},
	}

	apiCommands := &cli.Command{
		Name:        "api",
		Description: "Commands for managing the Harmonia API server.",
		Subcommands: []*cli.Command{
			{
				Name:        "start",
				Description: "Start the Harmonia API server",
				Action: func(c *cli.Context) error {
					log.Println("Starting Harmonia API server...")
					return nil
				},
			},
		},
	}

	app := &cli.App{
		Name:        "harmonia",
		Description: "Entrypoint of harmonia",
		Commands: []*cli.Command{
			cliCommands,
			apiCommands,
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
