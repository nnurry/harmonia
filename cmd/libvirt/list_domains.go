package libvirt

import (
	"bytes"
	"fmt"

	"github.com/nnurry/harmonia/internal/connection"
	"github.com/nnurry/harmonia/internal/service"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
	"libvirt.org/go/libvirt"
)

const LIST_LIBVIRT_DOMAIN_COMMAND = types.InternalCommandName("list Libvirt domains command")

type DomainMetadata struct {
	Name  string
	UUID  string
	State string
}

func NewDomainMetadata(domain libvirt.Domain) (*DomainMetadata, error) {
	domainName, err := domain.GetName()
	if err != nil {
		return nil, fmt.Errorf("fail to get name of domain: %v", err)
	}

	domainUuid, err := domain.GetUUIDString()
	if err != nil {
		return nil, fmt.Errorf("fail to get uuid of domain %v: %v", domainName, err)
	}

	domainState, reason, err := domain.GetState()
	if err != nil {
		return nil, fmt.Errorf("fail to get name of domain %v: %v (%v)", domainName, err, reason)
	}

	metadata := &DomainMetadata{
		Name: domainName,
		UUID: domainUuid,
	}
	switch domainState {
	case libvirt.DOMAIN_RUNNING:
		metadata.State = "running"
	case libvirt.DOMAIN_SHUTOFF:
		metadata.State = "shutoff"
	case libvirt.DOMAIN_SHUTDOWN:
		metadata.State = "shutdown"
	case libvirt.DOMAIN_CRASHED:
		metadata.State = "crashed"
	case libvirt.DOMAIN_NOSTATE:
		metadata.State = "nostate"
	case libvirt.DOMAIN_PAUSED:
		metadata.State = "paused"
	default:
		metadata.State = fmt.Sprintf("other (code=%v)", domainState)
	}
	return metadata, nil
}

func (metadata *DomainMetadata) ToString(pos int) string {
	buf := bytes.NewBufferString("")
	if pos > -1 {
		fmt.Fprintf(buf, "%v)	name: %v \n", pos, metadata.Name)
	} else {
		fmt.Fprintf(buf, "		name: %v \n", metadata.Name)
	}
	fmt.Fprintf(buf, "	uuid: %v\n", metadata.UUID)
	fmt.Fprintf(buf, "	state: %v\n", metadata.State)
	return buf.String()
}

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
		libvirtInternalConnection, ok := ctx.Context.Value(LIBVIRT_INTERNAL_CONNECTION_CTX_KEY).(*connection.Libvirt)
		if !ok {
			return fmt.Errorf("could not retrieve Libvirt internal connection from context")
		}

		libvirtService, err := service.NewLibvirt(libvirtInternalConnection)
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
			metadata, err := NewDomainMetadata(domain)
			if err != nil {
				fmt.Println("can't fetch metadata of domain -> ", err)
				continue
			}
			fmt.Println(metadata.ToString(i + 1))
		}

		return nil
	}
}

func (command *ListLibvirtDomainsCommand) Build() *cli.Command {
	return utils.ConvertInternalCommandToCliCommand(command)
}
