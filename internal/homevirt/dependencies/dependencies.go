package dependencies

import (
	"github.com/nnurry/harmonia/internal/homevirt/service"
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

func InitFromBuilder(builder *DependenciesBuilder) (*Dependencies, error) {
	libvirtSvc, err := service.NewLibvirt(builder.conn)
	if err != nil {
		return nil, err
	}
	return &Dependencies{
		libvirtSvc: libvirtSvc,
	}, nil
}
