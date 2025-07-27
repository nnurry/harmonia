package libvirt

import (
	"bytes"
	"context"
	"fmt"

	"github.com/nnurry/harmonia/internal/connection"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	LIBVIRT_COMMAND = types.InternalCommandName("Libvirt command")
)

const (
	LIBVIRT_INTERNAL_CONNECTION_CTX_KEY = types.InternalCommandCtxKey("libvirtInternalConnection")
)

type LibvirtCommand struct {
	config connection.LibvirtConfig
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
			Destination: &command.config.ConnectionUrl,
		},
		&cli.StringFlag{
			Name:        "keyfile-path",
			Destination: &command.config.KeyfilePath,
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
		buf := bytes.NewBufferString("")
		for _, subcmd := range ctx.Command.Subcommands {
			fmt.Fprintf(buf, "- %v\n", subcmd.Name)
		}
		return fmt.Errorf("use subcommands instead:\n%v", buf.String())
	}
}

func (command *LibvirtCommand) Build() *cli.Command {
	cliCommand := utils.ConvertInternalCommandToCliCommand(command)
	cliCommand.Before = func(ctx *cli.Context) error {
		libvirtConnection, err := connection.NewLibvirt(command.config)

		if err != nil {
			return fmt.Errorf("could not establish Libvirt connection: %v", err)
		}

		log.Info().Msgf("Setting Libvirt connection with URL = %v\n", libvirtConnection.URL())

		ctx.Context = context.WithValue(ctx.Context, LIBVIRT_INTERNAL_CONNECTION_CTX_KEY, libvirtConnection)
		return nil

	}

	return cliCommand
}
