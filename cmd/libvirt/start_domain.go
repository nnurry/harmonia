package libvirt

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/builder"
	"github.com/nnurry/harmonia/internal/service"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

type StartLibvirtDomainCommand struct {
}

func (command *StartLibvirtDomainCommand) Description() string {
	return "Start a domain in Libvirt"
}

func (command *StartLibvirtDomainCommand) Signature() string {
	return "start-domain"
}

func (command *StartLibvirtDomainCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (command *StartLibvirtDomainCommand) Subcommands() []*cli.Command {
	return []*cli.Command{}
}

func (command *StartLibvirtDomainCommand) Handler() func(ctx *cli.Context) error {
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

		err = libvirtService.StartDomainWithName(domainName)
		if err != nil {
			return fmt.Errorf("could not start domain %v: %v", domainName, err)
		}

		return nil
	}
}

func (command *StartLibvirtDomainCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
