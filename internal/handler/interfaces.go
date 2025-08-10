package handler

import (
	"context"
	"io"

	"github.com/nnurry/harmonia/internal/builder"
	"github.com/nnurry/harmonia/internal/service/cloudinit"
	"libvirt.org/go/libvirt"
)

type LibvirtService interface {
	GetDomainByName(name string) (*libvirt.Domain, error)
	DefineDomainFromBuilder(domainBuilder *builder.LibvirtDomainBuilder) (*libvirt.Domain, error)
}

type CloudInitService interface {
	SetUserData(userData cloudinit.UserData)
	SetMetaData(userData cloudinit.MetaData)
	SetNetworkConfig(userData cloudinit.NetworkConfig)
	WriteToDisk(ctx context.Context, basePath string, filename string) (string, error)
	RemoveFromDisk(ctx context.Context, basePath string) error
}

type ShellProcessor interface {
	Name() string
	Execute(ctx context.Context, stdout io.Writer, stderr io.Writer, command string, arguments ...string) error
}
