package libvirt

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/homevirt/builder"
	"github.com/nnurry/harmonia/internal/homevirt/service"
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
			Name:     "ci-iso-path",
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
		if ctx.NArg() < 1 {
			return fmt.Errorf("missing <base domain name> (need 1 base domain to build other child domains)")
		}

		baseDomainName := ctx.Args().First()
		if baseDomainName == "" {
			return fmt.Errorf("<base domain name> is empty")
		}

		connectBuilder, ok := ctx.Context.Value(LIBVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY).(*builder.LibvirtConnectBuilder)
		if !ok {
			return fmt.Errorf("could not retrieve Libvirt connect builder from context")
		}

		libvirtService, err := service.NewLibvirtFromConnectBuilder(connectBuilder)
		if err != nil {
			return err
		}

		baseDomain, err := libvirtService.GetDomainByName(baseDomainName)
		if err != nil {
			return err
		}

		domainBuilder, err := builder.NewLibvirtDomainBuilder(
			baseDomain,
			[]*builder.DomainBuilderFlag{
				builder.SET_VM_NAME,
				builder.SET_CI_DISK_PATH,
				builder.SET_QCOW2_DISK_PATH,
			},
			false,
		)

		if ctx.IsSet("name") {
			domainBuilder = domainBuilder.WithDomainName(ctx.String("name"))
		}

		if ctx.IsSet("cpus") {
			domainBuilder = domainBuilder.WithNumOfCpus(ctx.Int("cpus"))
		}

		if ctx.IsSet("memory") {
			domainBuilder = domainBuilder.WithMemory(uint(ctx.Float64("memory")*1024*1024), "KiB")
		}

		if ctx.IsSet("qcow2-disk-path") {
			domainBuilder = domainBuilder.WithQcow2DiskPath(ctx.String("qcow2-disk-path"))
		}

		if ctx.IsSet("ci-iso-path") {
			domainBuilder = domainBuilder.WithCiDiskPath(ctx.String("ci-iso-path"))
		}

		if err != nil {
			return err
		}

		_, err = libvirtService.DefineDomainFromBuilder(domainBuilder)

		if err != nil {
			return err
		}

		return nil
	}
}

func (command *DefineLibvirtDomainCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
