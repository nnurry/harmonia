package libvirt

import (
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

const DEFINE_LIBVIRT_DOMAIN_COMMAND = types.InternalCommandName("define Libvirt domain command")

type DefineLibvirtDomainCommand struct{}

func (command *DefineLibvirtDomainCommand) Description() string {
	return "Define a domain in Libvirt"
}

func (command *DefineLibvirtDomainCommand) Signature() string {
	return "define-domain"
}

func (command *DefineLibvirtDomainCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return nil
	}
}

func (command *DefineLibvirtDomainCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (command *DefineLibvirtDomainCommand) Subcommands() []*cli.Command {
	return []*cli.Command{}
}

func (command *DefineLibvirtDomainCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
