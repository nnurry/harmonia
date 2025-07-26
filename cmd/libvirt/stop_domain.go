package libvirt

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/homevirt/builder"
	"github.com/nnurry/harmonia/internal/homevirt/service"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

type StopLibvirtDomainCommand struct {
}

func (command *StopLibvirtDomainCommand) Description() string {
	return "Shutdown a domain in Libvirt"
}

func (command *StopLibvirtDomainCommand) Signature() string {
	return "stop-domain"
}

func (command *StopLibvirtDomainCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (command *StopLibvirtDomainCommand) Subcommands() []*cli.Command {
	return []*cli.Command{}
}

func (command *StopLibvirtDomainCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if ctx.NArg() < 1 {
			return fmt.Errorf("missing <domain name>")
		}

		domainName := ctx.Args().First()
		if domainName == "" {
			return fmt.Errorf("<domain name> is empty")
		}

		connectBuilder, ok := ctx.Context.Value(LIBVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY).(*builder.LibvirtConnectBuilder)
		if !ok {
			return fmt.Errorf("could not retrieve Libvirt connect builder from context")
		}

		libvirtService, err := service.NewLibvirtFromConnectBuilder(connectBuilder)
		if err != nil {
			return err
		}
		defer libvirtService.Cleanup()

		err = libvirtService.StopDomainWithName(domainName)
		if err != nil {
			return fmt.Errorf("could not start domain %v: %v", domainName, err)
		}

		return nil
	}
}

func (command *StopLibvirtDomainCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
