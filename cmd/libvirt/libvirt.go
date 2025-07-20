package libvirt

import (
	"context"
	"fmt"

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
	hypervisor    string
	transportType string
	user          string
	host          string
	keyfilePath   string

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
			Name:        "hypervisor",
			DefaultText: builder.LIBVIRT_CONNECT_URL_DEFAULT_HYPERVISOR,
			Usage:       "Set type of hypervisor. Omit to use 'qemu' (no support for xen though, hehe)",
			Destination: &command.hypervisor,
		},
		&cli.StringFlag{
			Name:        "transport-type",
			Usage:       "Use 'ssh' to connect remotely. I do not support other type of connections. If this is run on the same machine as hypervisor, omitting should be fine",
			Destination: &command.transportType,
		},
		&cli.StringFlag{
			Name:        "user",
			Usage:       "Set system user. Omit to use default",
			Destination: &command.user,
		},
		&cli.StringFlag{
			Name:        "host",
			Usage:       "Host address of hypervisor. Omit to use 'localhost' (only when hypervisor is on the same machine, otherwise recommend to fill the address)",
			Destination: &command.host,
		},
		&cli.StringFlag{
			Name:        "keyfile-path",
			Usage:       "Path to SSH private key for 'ssh' transport type. Omit if SSH config already has matching identity",
			Destination: &command.keyfilePath,
		},
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
		if command.connectUrl != "" {
			connectBuilder, err := builder.NewLibvirtConnectBuilderFromConnectUrl(command.connectUrl)

			if err != nil {
				return err
			}

			ctx.Context = context.WithValue(ctx.Context, LIBVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY, connectBuilder)
			return nil
		}

		connectBuilder, err := builder.NewLibvirtConnectBuilder(
			[]*builder.ConnectUrlBuilderFlag{
				builder.SET_TRANSPORT_CONF,
				builder.SET_USER,
				builder.SET_HOST,
				builder.SET_KEYFILE,
			},
			false,
		)

		if err != nil {
			return err
		}

		if command.hypervisor != "" {
			connectBuilder = connectBuilder.WithHypervisor(command.hypervisor)
		}

		if command.transportType != "" {
			connectBuilder = connectBuilder.WithTransportType(command.transportType)
		}

		if command.user != "" {
			connectBuilder = connectBuilder.WithUser(command.user)
		}

		if command.host != "" {
			connectBuilder = connectBuilder.WithHost(command.host)
		}

		if command.keyfilePath != "" {
			connectBuilder = connectBuilder.WithKeyfilePath(command.keyfilePath)
		}

		if err = connectBuilder.Verify(); err != nil {
			return err
		}

		ctx.Context = context.WithValue(ctx.Context, LIBVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY, connectBuilder)
		return nil
	}

	return cliCommand
}
