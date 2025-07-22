package homevirt

import (
	"fmt"

	"github.com/nnurry/harmonia/internal/homevirt/builder"
	"github.com/nnurry/harmonia/internal/homevirt/service"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

const LIST_LIBVIRT_DOMAIN_COMMAND = types.InternalCommandName("list Libvirt domains command")

type ListLibvirtDomainsCommand struct {
	isListAll bool
}

func (command *ListLibvirtDomainsCommand) Description() string {
	return "List domains in Libvirt"
}

func (command *ListLibvirtDomainsCommand) Signature() string {
	return "list-domains"
}

func (command *ListLibvirtDomainsCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "all",
			Usage:       "This will include inactive domains as well.",
			Destination: &command.isListAll,
		},
	}
}

func (command *ListLibvirtDomainsCommand) Subcommands() []*cli.Command {
	return []*cli.Command{}
}

func (command *ListLibvirtDomainsCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		connectBuilder, ok := ctx.Context.Value(HOMEVIRT_COMMAND_CONNECT_BUILDER_CTX_KEY).(*builder.LibvirtConnectBuilder)
		if !ok {
			return fmt.Errorf("could not retrieve Libvirt connect builder from context")
		}

		libvirtService, err := service.NewLibvirtFromConnectBuilder(connectBuilder)
		if err != nil {
			return err
		}
		defer libvirtService.Cleanup()
		domains, err := libvirtService.ListDomains(command.isListAll)

		if err != nil {
			return fmt.Errorf("could not list domains: %v", err)
		}

		fmt.Println("List of domains:")

		for i, domain := range domains {
			domainName, _ := domain.GetName()
			domainUuid, _ := domain.GetUUIDString()
			fmt.Printf("%v)	name: %v\n", i+1, domainName)
			fmt.Printf("	uuid: %v\n", domainUuid)
		}

		return nil
	}
}

func (command *ListLibvirtDomainsCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
