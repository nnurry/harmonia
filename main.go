package main

import (
	"log"
	"os"

	"github.com/nnurry/harmonia/cmd"
	libvirtcmd "github.com/nnurry/harmonia/cmd/libvirt"
	shellcmd "github.com/nnurry/harmonia/cmd/shell"
	"github.com/urfave/cli/v2"
)

func main() {
	cliCommands := &cli.Command{
		Name:        "cli",
		Description: "Commands for interacting with Harmonia's features directly.",
		Subcommands: []*cli.Command{
			cmd.GetCliCommand(libvirtcmd.LIBVIRT_COMMAND),
			cmd.GetCliCommand(shellcmd.SHELL_COMMAND),
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
		Name:                 "harmonia",
		Description:          "Entrypoint of harmonia",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			cliCommands,
			apiCommands,
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
