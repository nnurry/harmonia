package libvirt

import (
	"context"
	"fmt"
	"log"

	"github.com/nnurry/harmonia/internal/homevirt/builder"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

const (
	LIBVIRT_COMMAND = types.InternalCommandName("Libvirt command")
)

const (
	LIBVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY = types.InternalCommandCtxKey("connectBuilder")
)

type LibvirtCommand struct {
	connectUrl string
}

func (command *LibvirtCommand) Description() string {
	return "Libvirt command entrypoint"
}

func (command *LibvirtCommand) Signature() string {
	return "libvirt"
}

func (command *LibvirtCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "connect-url",
			Usage:       "Set Libvirt connect URL. Using this flag would discard all other connect flags",
			Destination: &command.connectUrl,
		},
	}
}

func (command *LibvirtCommand) Subcommands() []*cli.Command {
	return []*cli.Command{
		(&DefineLibvirtDomainCommand{}).Build(),
		(&RemoveLibvirtDomainCommand{}).Build(),
		(&ListLibvirtDomainsCommand{}).Build(),
		(&StartLibvirtDomainCommand{}).Build(),
		(&StopLibvirtDomainCommand{}).Build(),
	}
}

func (command *LibvirtCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return fmt.Errorf("use subcommands instead")
	}
}

func (command *LibvirtCommand) Build() *cli.Command {
	cliCommand := utils.ConvertInternalCommandToCliCommand(command)
	cliCommand.Before = func(ctx *cli.Context) error {
		connectBuilder, err := builder.NewLibvirtConnectBuilderFromConnectUrl(command.connectUrl)

		if err != nil {
			return err
		}

		log.Printf("Setting Libvirt connection with URL = %v\n", command.connectUrl)

		ctx.Context = context.WithValue(ctx.Context, LIBVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY, connectBuilder)
		return nil

	}

	return cliCommand
}
