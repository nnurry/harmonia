package service

import (
	"fmt"
	"net/url"
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

type LibvirtConnUrlBuilder struct {
	transportType string
	hypervisor    string
	host          string
	user          string
	path          string
	keyfilePath   string

	requiredFlags map[ConnUrlBuilderFlag]bool
}

func NewLibvirtConnUrlBuilder() *LibvirtConnUrlBuilder {
	builder := &LibvirtConnUrlBuilder{
		hypervisor:    "qemu",
		path:          "system",
		host:          "localhost",
		requiredFlags: map[ConnUrlBuilderFlag]bool{},
	}
	return builder
}

func (builder *LibvirtConnUrlBuilder) WithTransportType(transportType string) *LibvirtConnUrlBuilder {
	builder.transportType = transportType
	return builder
}

func (builder *LibvirtConnUrlBuilder) WithHypervisor(hypervisor string) *LibvirtConnUrlBuilder {
	builder.hypervisor = hypervisor
	return builder
}

func (builder *LibvirtConnUrlBuilder) WithUser(user string) *LibvirtConnUrlBuilder {
	builder.user = user
	return builder
}

func (builder *LibvirtConnUrlBuilder) WithHost(host string) *LibvirtConnUrlBuilder {
	builder.host = host
	return builder
}

func (builder *LibvirtConnUrlBuilder) WithConnectAsRoot(connectAsRoot bool) *LibvirtConnUrlBuilder {
	if connectAsRoot {
		builder.path = "system"
	} else {
		builder.path = "domain"
	}
	return builder
}

func (builder *LibvirtConnUrlBuilder) WithKeyfilePath(keyfilePath string) *LibvirtConnUrlBuilder {
	builder.keyfilePath = keyfilePath
	return builder
}

func (builder *LibvirtConnUrlBuilder) BuildConnStr() (string, error) {
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
