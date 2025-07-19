package builder

import (
	// "errors"

	"net/url"

	"github.com/nnurry/harmonia/pkg/types"
	"libvirt.org/go/libvirt"
)

const (
	LIBVIRT_CONNECT_URL_DEFAULT_HYPERVISOR = "qemu"
	LIBVIRT_CONNECT_URL_DEFAULT_PATH       = "system"
	LIBVIRT_CONNECT_URL_DEFAULT_HOST       = "localhost"
)

type ConnectUrlBuilderFlag struct {
	name string
}

func (flag *ConnectUrlBuilderFlag) Name() string {
	return flag.name
}

var (
	SET_HYPERVISOR     = &ConnectUrlBuilderFlag{name: "set hypervisor"}
	SET_TRANSPORT_CONF = &ConnectUrlBuilderFlag{name: "set transport configuration"}
	SET_HOST           = &ConnectUrlBuilderFlag{name: "set host"}
	SET_USER           = &ConnectUrlBuilderFlag{name: "set user"}
	SET_ROOT_CONNECT   = &ConnectUrlBuilderFlag{name: "set root connect"}
	SET_KEYFILE        = &ConnectUrlBuilderFlag{name: "set keyfile"}
)

type LibvirtConnectBuilder struct {
	transportType string
	hypervisor    string
	host          string
	user          string
	path          string
	keyfilePath   string

	builderFlagMap *types.BuilderFlagMap
}

func NewLibvirtConnectBuilder(useDefaultBuilderFlags bool, requiredFlags []*ConnectUrlBuilderFlag) (*LibvirtConnectBuilder, error) {
	castedRequiredFlags := []types.BuilderFlag{}

	for _, flag := range requiredFlags {
		castedRequiredFlags = append(castedRequiredFlags, types.BuilderFlag(flag))
	}

	builder := &LibvirtConnectBuilder{
		hypervisor: LIBVIRT_CONNECT_URL_DEFAULT_HYPERVISOR,
		path:       LIBVIRT_CONNECT_URL_DEFAULT_PATH,
		host:       LIBVIRT_CONNECT_URL_DEFAULT_HOST,
	}
	builderFlagMap, err := types.NewFlagMapFromBuilderFlags(
		castedRequiredFlags,
		builder.getDefaultBuilderFlags(),
		useDefaultBuilderFlags,
	)

	if err != nil {
		return nil, err
	}

	builder.builderFlagMap = builderFlagMap
	return builder, nil
}

func (builder *LibvirtConnectBuilder) getDefaultBuilderFlags() []types.BuilderFlag {
	return []types.BuilderFlag{}
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

func (builder *LibvirtConnectBuilder) Verify() error {
	return builder.builderFlagMap.Verify()
}

func (builder *LibvirtConnectBuilder) BuildConnectURL() (string, error) {
	if err := builder.Verify(); err != nil {
		return "", err
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
