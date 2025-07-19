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

func (command *DefineLibvirtDomainCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Required: true,
			Usage:    "Set name of new domain. Required, sadly",
		},
		&cli.IntFlag{
			Name:  "cpus",
			Value: 2,
			Usage: "Set number of CPUs. Default to 2 CPUs",
		},
		&cli.Float64Flag{
			Name:  "memory",
			Value: 8,
			Usage: "Set RAM amount in GiB (accept float value). Default to 8GiB",
		},
		&cli.StringFlag{
			Name:     "qcow2-disk-path",
			Required: true,
			Usage:    "Set path to qcow2 disk file (where domain stores data). It is created outside of this scope via qemu-img (COW) or cp (deepcopy). Required, sadly",
		},
		&cli.StringFlag{
			Name:     "keyfile-path",
			Required: true,
			Usage:    "Set path to cloud-init ISO file. It is created outside of this scope via mkisofs. Required, sadly",
		},
	}
}

func (command *DefineLibvirtDomainCommand) Subcommands() []*cli.Command {
	return []*cli.Command{}
}

func (command *DefineLibvirtDomainCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return nil
	}
}

func (command *DefineLibvirtDomainCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
