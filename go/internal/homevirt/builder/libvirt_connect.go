package builder

import (
	"fmt"
	"net/url"

	"github.com/nnurry/harmonia/pkg/types"
	"libvirt.org/go/libvirt"
)

type ConnectUrlBuilderFlag types.BuilderFlag

const (
	SET_HYPERVISOR     = ConnectUrlBuilderFlag("set hypervisor")
	SET_TRANSPORT_CONF = ConnectUrlBuilderFlag("set transport configuration")
	SET_HOST           = ConnectUrlBuilderFlag("set host")
	SET_USER           = ConnectUrlBuilderFlag("set user")
	SET_ROOT_CONNECT   = ConnectUrlBuilderFlag("set root connect")
	SET_KEYFILE        = ConnectUrlBuilderFlag("set keyfile")
)

type LibvirtConnectBuilder struct {
	transportType string
	hypervisor    string
	host          string
	user          string
	path          string
	keyfilePath   string

	requiredFlagMap map[ConnectUrlBuilderFlag]bool
}

func NewLibvirtConnectBuilder(requiredFlags ...ConnectUrlBuilderFlag) *LibvirtConnectBuilder {
	if len(requiredFlags) < 1 {
		requiredFlags = []ConnectUrlBuilderFlag{}
	}

	requiredFlagMap := make(map[ConnectUrlBuilderFlag]bool)
	for _, flag := range requiredFlags {
		requiredFlagMap[flag] = false
	}

	builder := &LibvirtConnectBuilder{
		hypervisor:      "qemu",
		path:            "system",
		host:            "localhost",
		requiredFlagMap: requiredFlagMap,
	}
	return builder
}

func (builder *LibvirtConnectBuilder) WithTransportType(transportType string) *LibvirtConnectBuilder {
	builder.transportType = transportType
	return builder
}

func (builder *LibvirtConnectBuilder) WithHypervisor(hypervisor string) *LibvirtConnectBuilder {
	builder.hypervisor = hypervisor
	return builder
}

func (builder *LibvirtConnectBuilder) WithUser(user string) *LibvirtConnectBuilder {
	builder.user = user
	return builder
}

func (builder *LibvirtConnectBuilder) WithHost(host string) *LibvirtConnectBuilder {
	builder.host = host
	return builder
}

func (builder *LibvirtConnectBuilder) WithConnectAsRoot(connectAsRoot bool) *LibvirtConnectBuilder {
	if connectAsRoot {
		builder.path = "system"
	} else {
		builder.path = "domain"
	}
	return builder
}

func (builder *LibvirtConnectBuilder) WithKeyfilePath(keyfilePath string) *LibvirtConnectBuilder {
	builder.keyfilePath = keyfilePath
	return builder
}

func (builder *LibvirtConnectBuilder) BuildConnectURL() (string, error) {
	unsatisfiedFlags := []ConnectUrlBuilderFlag{}
	for flag, satisfied := range builder.requiredFlagMap {
		if !satisfied {
			unsatisfiedFlags = append(unsatisfiedFlags, flag)
		}
	}

	if len(unsatisfiedFlags) > 0 {
		return "", fmt.Errorf("flag [%v] not satisfied", unsatisfiedFlags)
	}

	scheme := builder.hypervisor
	if builder.transportType != "" {
		scheme = scheme + "+" + builder.transportType
	}

	host := builder.host
	if host != "/" {
		if builder.user != "" {
			host = builder.user + "@" + builder.host
		}
	}

	queryStr := url.Values{}
	if builder.keyfilePath != "" {
		queryStr.Add("keyfile", builder.keyfilePath)
	}

	connUrl := url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     builder.path,
		RawQuery: queryStr.Encode(),
	}

	return connUrl.String(), nil
}

func (builder *LibvirtConnectBuilder) Build() (*libvirt.Connect, error) {
	connectUrl, err := builder.BuildConnectURL()
	if err != nil {
		return nil, err
	}

	connect, err := libvirt.NewConnect(connectUrl)
	if err != nil {
		return nil, err
	}

	return connect, nil
}
