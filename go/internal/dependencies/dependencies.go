package dependencies

import (
	"github.com/nnurry/harmonia/internal/service"
	"libvirt.org/go/libvirt"
)

type DependenciesBuilder struct {
	conn *libvirt.Connect
}

type Dependencies struct {
	libvirtSvc *service.Libvirt
}

func InitBuilder(conn *libvirt.Connect) *DependenciesBuilder {
	return &DependenciesBuilder{conn: conn}
}

func InitFromBuilder(builder *DependenciesBuilder) *Dependencies {
	libvirtSvc := service.NewLibvirt(builder.conn)
	return &Dependencies{
		libvirtSvc: libvirtSvc,
	}
}
