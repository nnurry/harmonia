package service

import (
	"fmt"
	"net/url"

	"libvirt.org/go/libvirt"
)

type ConnUrlBuilderFlag string

const (
	SET_HYPERVISOR     = ConnUrlBuilderFlag("set hypervisor")
	SET_TRANSPORT_CONF = ConnUrlBuilderFlag("set transport configuration")
	SET_HOST           = ConnUrlBuilderFlag("set host")
	SET_USER           = ConnUrlBuilderFlag("set user")
	SET_ROOT_CONNECT   = ConnUrlBuilderFlag("set root connect")
	SET_KEYFILE        = ConnUrlBuilderFlag("set keyfile")
)

type LibvirtConn struct {
}

type LibvirtConnectBuilder struct {
	transportType string
	hypervisor    string
	host          string
	user          string
	path          string
	keyfilePath   string

	requiredFlags map[ConnUrlBuilderFlag]bool
}

func NewLibvirtConnectBuilder() *LibvirtConnectBuilder {
	builder := &LibvirtConnectBuilder{
		hypervisor:    "qemu",
		path:          "system",
		host:          "localhost",
		requiredFlags: map[ConnUrlBuilderFlag]bool{},
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
	unsatisfiedFlags := []ConnUrlBuilderFlag{}
	for flag, satisfied := range builder.requiredFlags {
		if !satisfied {
			unsatisfiedFlags = append(unsatisfiedFlags, flag)
		}
	}

	if len(unsatisfiedFlags) > 0 {
		return "", fmt.Errorf("failed to build Libvirt conn URL: flag [%v] not satisfied", unsatisfiedFlags)
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
