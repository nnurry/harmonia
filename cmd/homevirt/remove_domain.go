package homevirt

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/homevirt/builder"
	"github.com/nnurry/harmonia/internal/homevirt/service"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

type RemoveLibvirtDomainCommand struct {
}

func (command *RemoveLibvirtDomainCommand) Description() string {
	return "Remove a domain in Libvirt"
}

func (command *RemoveLibvirtDomainCommand) Signature() string {
	return "remove-domain"
}

func (command *RemoveLibvirtDomainCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (command *RemoveLibvirtDomainCommand) Subcommands() []*cli.Command {
	return []*cli.Command{}
}

func (command *RemoveLibvirtDomainCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if ctx.NArg() < 1 {
			return fmt.Errorf("missing <domain name>")
		}

		domainName := ctx.Args().First()
		if domainName == "" {
			return fmt.Errorf("<domain name> is empty")
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

		domain, err := libvirtService.GetDomainByName(domainName)

		if err != nil {
			return err
		}

		err = domain.Undefine()
		if err != nil {
			return fmt.Errorf("could not remove domain %v: %v", domainName, err)
		}

		return nil
	}
}

func (command *RemoveLibvirtDomainCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
