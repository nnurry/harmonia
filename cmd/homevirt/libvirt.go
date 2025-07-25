package homevirt

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
	HOMEVIRT_COMMAND = types.InternalCommandName("Homevirt command")
)

const (
	HOMEVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY = types.InternalCommandCtxKey("connectBuilder")
)

type HomevirtCommand struct {
	connectUrl string
}

func (command *HomevirtCommand) Description() string {
	return "Homevirt command entrypoint"
}

func (command *HomevirtCommand) Signature() string {
	return "homevirt"
}

func (command *HomevirtCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "connect-url",
			Usage:       "Set Homevirt connect URL. Using this flag would discard all other connect flags",
			Destination: &command.connectUrl,
		},
	}
}

func (command *HomevirtCommand) Subcommands() []*cli.Command {
	return []*cli.Command{
		(&DefineLibvirtDomainCommand{}).Build(),
		(&RemoveLibvirtDomainCommand{}).Build(),
		(&ListLibvirtDomainsCommand{}).Build(),
		(&StartLibvirtDomainCommand{}).Build(),
		(&StopLibvirtDomainCommand{}).Build(),
	}
}

func (command *HomevirtCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return fmt.Errorf("use subcommands instead")
	}
}

func (command *HomevirtCommand) Build() *cli.Command {
	cliCommand := utils.ConvertInternalCommandToCliCommand(command)
	cliCommand.Before = func(ctx *cli.Context) error {
		connectBuilder, err := builder.NewLibvirtConnectBuilderFromConnectUrl(command.connectUrl)

		if err != nil {
			return err
		}

		log.Printf("Setting Libvirt connection with URL = %v\n", command.connectUrl)

		ctx.Context = context.WithValue(ctx.Context, HOMEVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY, connectBuilder)
		return nil

	}

	return cliCommand
}
