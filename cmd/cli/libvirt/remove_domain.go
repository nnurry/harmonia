package libvirt

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/connection"
	"github.com/nnurry/harmonia/internal/service"
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

		libvirtInternalConnection, ok := ctx.Context.Value(LIBVIRT_INTERNAL_CONNECTION_CTX_KEY).(*connection.Libvirt)
		if !ok {
			return fmt.Errorf("could not retrieve Libvirt internal connection from context")
		}

		libvirtService, err := service.NewLibvirt(libvirtInternalConnection)
		if err != nil {
			return err
		}
		defer libvirtService.Cleanup()

		err = libvirtService.RemoveDomainByName(domainName)
		if err != nil {
			return fmt.Errorf("could not remove domain %v: %v", domainName, err)
		}

		return nil
	}
}

func (command *RemoveLibvirtDomainCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
