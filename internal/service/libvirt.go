package service

import (
	"github.com/nnurry/harmonia/internal/builder"
	"github.com/nnurry/harmonia/internal/connection"
	"libvirt.org/go/libvirt"
)

const (
	DEFAULT_LIBVIRT_QEMU_DISK_BASE_PATH = "/var/lib/libvirt/images"
)

type Libvirt struct {
	*connection.Libvirt
}

func NewLibvirt(connection *connection.Libvirt) (*Libvirt, error) {
	return &Libvirt{connection}, nil
}

func (service *Libvirt) GetDomainByName(name string) (*libvirt.Domain, error) {
	return service.Connect().LookupDomainByName(name)
}

func (service *Libvirt) DefineDomainFromXMLString(xmlString string) (*libvirt.Domain, error) {
	return service.Connect().DomainDefineXML(xmlString)
}

func (service *Libvirt) DefineDomainFromBuilder(domainBuilder *builder.LibvirtDomainBuilder) (*libvirt.Domain, error) {
	return domainBuilder.Build(service.Connect())
}

func (service *Libvirt) RemoveDomainByName(name string) error {
	domain, err := service.GetDomainByName(name)

	if err != nil {
		return err
	}

	return domain.Undefine()
}

func (service *Libvirt) StartDomainWithName(name string) error {
	domain, err := service.GetDomainByName(name)

	if err != nil {
		return err
	}

	return domain.Create()
}

func (service *Libvirt) StopDomainWithName(name string) error {
	domain, err := service.GetDomainByName(name)

	if err != nil {
		return err
	}

	return domain.Shutdown()
}

func (service *Libvirt) ListDomains(includeInactive bool) ([]libvirt.Domain, error) {
	flags := libvirt.CONNECT_LIST_DOMAINS_ACTIVE
	if includeInactive {
		flags |= libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	}
	return service.Connect().ListAllDomains(flags)
}
