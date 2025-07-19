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

type LibvirtCommand struct{}

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
		},
		&cli.StringFlag{
			Name:  "transport-type",
			Usage: "Use 'ssh' to connect remotely. I do not support other type of connections. If this is run on the same machine as hypervisor, omitting should be fine",
		},
		&cli.StringFlag{
			Name:  "user",
			Usage: "Set system user. Omit to use default",
		},
		&cli.StringFlag{
			Name:  "host",
			Usage: "Host address of hypervisor. Omit to use 'localhost' (only when hypervisor is on the same machine, otherwise recommend to fill the address)",
		},
		&cli.StringFlag{
			Name:  "keyfile-path",
			Usage: "Path to SSH private key for 'ssh' transport type. Omit if SSH config already has matching identity",
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
		connectBuilder, err := builder.NewLibvirtConnectBuilder(
			false,
			[]*builder.ConnectUrlBuilderFlag{
				builder.SET_TRANSPORT_CONF,
				builder.SET_USER,
				builder.SET_HOST,
				builder.SET_KEYFILE,
			},
		)

		if err != nil {
			return err
		}

		if ctx.IsSet("hypervisor") {
			connectBuilder = connectBuilder.WithHypervisor(ctx.String("hypervisor"))
		}

		if ctx.IsSet("transport-type") {
			connectBuilder = connectBuilder.WithTransportType(ctx.String("transport-type"))
		}

		if ctx.IsSet("user") {
			connectBuilder = connectBuilder.WithUser(ctx.String("user"))
		}

		if ctx.IsSet("host") {
			connectBuilder = connectBuilder.WithHost(ctx.String("host"))
		}

		if ctx.IsSet("keyfile-path") {
			connectBuilder = connectBuilder.WithKeyfilePath(ctx.String("keyfile-path"))
		}

		if err = connectBuilder.Verify(); err != nil {
			return nil
		}

		ctx.Context = context.WithValue(ctx.Context, LIBVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY, connectBuilder)
		return nil
	}

	return cliCommand
}
