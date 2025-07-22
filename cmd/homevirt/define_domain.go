package homevirt

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/homevirt/builder"
	"github.com/nnurry/harmonia/internal/homevirt/service"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

const DEFINE_LIBVIRT_DOMAIN_COMMAND = types.InternalCommandName("define Libvirt domain command")

type DefineLibvirtDomainCommand struct {
	domainName    string
	cpus          int
	memory        float64
	qcow2DiskPath string
	ciIsoPath     string
}

func (command *DefineLibvirtDomainCommand) Description() string {
	return "Define a domain in Libvirt"
}

func (command *DefineLibvirtDomainCommand) Signature() string {
	return "define-domain"
}

func (command *DefineLibvirtDomainCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "name",
			Required:    true,
			Usage:       "Set name of new domain. Required, sadly",
			Destination: &command.domainName,
		},
		&cli.IntFlag{
			Name:        "cpus",
			Value:       2,
			Usage:       "Set number of CPUs. Default to 2 CPUs",
			Destination: &command.cpus,
		},
		&cli.Float64Flag{
			Name:        "memory",
			Value:       8,
			Usage:       "Set RAM amount in GiB (accept float value). Default to 8GiB",
			Destination: &command.memory,
		},
		&cli.StringFlag{
			Name:        "qcow2-disk-path",
			Usage:       "Set path to qcow2 disk file (where domain stores data). It is created outside of this scope via qemu-img (COW) or cp (deepcopy). Default to '/var/lib/libvirt/images/<new domain name>.qcow2'",
			Destination: &command.qcow2DiskPath,
		},
		&cli.StringFlag{
			Name:        "ci-iso-path",
			Usage:       "Set path to cloud-init ISO file. It is created outside of this scope via mkisofs. Default to '/var/lib/libvirt/images/<new domain name>.iso'",
			Destination: &command.ciIsoPath,
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

		connectBuilder, ok := ctx.Context.Value(HOMEVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY).(*builder.LibvirtConnectBuilder)
		if !ok {
			return fmt.Errorf("could not retrieve Libvirt connect builder from context")
		}

		libvirtService, err := service.NewLibvirtFromConnectBuilder(connectBuilder)
		if err != nil {
			return err
		}
		defer libvirtService.Cleanup()

		baseDomain, err := libvirtService.GetDomainByName(baseDomainName)
		if err != nil {
			return err
		}

		domainBuilder, err := builder.NewLibvirtDomainBuilder(
			baseDomain,
			[]*builder.DomainBuilderFlag{
				builder.SET_VM_NAME,
			},
			false,
		)

		domainBuilder = domainBuilder.WithDomainName(command.domainName)

		if command.cpus != 0 {
			domainBuilder = domainBuilder.WithNumOfCpus(command.cpus)
		}

		if command.memory != 0.0 {
			domainBuilder = domainBuilder.WithMemory(uint(command.memory*1024*1024), "KiB")
		}

		if command.qcow2DiskPath == "" {
			command.qcow2DiskPath = fmt.Sprintf(
				"%s/%s.qcow2",
				service.DEFAULT_LIBVIRT_QEMU_DISK_BASE_PATH,
				command.domainName,
			)

		}
		domainBuilder = domainBuilder.WithQcow2DiskPath(command.qcow2DiskPath)

		if command.ciIsoPath == "" {
			command.ciIsoPath = fmt.Sprintf(
				"%s/%s-ci-data.iso",
				service.DEFAULT_LIBVIRT_QEMU_DISK_BASE_PATH,
				command.domainName,
			)
		}
		domainBuilder = domainBuilder.WithCiDiskPath(command.ciIsoPath)

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
