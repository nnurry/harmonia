package libvirt

import (
	"fmt"

	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

const LIBVIRT_COMMAND = types.InternalCommandName("Libvirt command")

type LibvirtCommand struct{}

func (command *LibvirtCommand) Description() string {
	return "Libvirt command entrypoint"
}

func (command *LibvirtCommand) Signature() string {
	return "libvirt"
}

func (command *LibvirtCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return fmt.Errorf("use subcommands instead")
	}
}

func (command *LibvirtCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (command *LibvirtCommand) Subcommands() []*cli.Command {
	return []*cli.Command{
		(&DefineLibvirtDomainCommand{}).Build(),
	}
}
