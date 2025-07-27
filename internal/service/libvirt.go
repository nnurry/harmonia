package service

import (
	"github.com/nnurry/harmonia/internal/builder"
	"libvirt.org/go/libvirt"
)

const (
	DEFAULT_LIBVIRT_QEMU_DISK_BASE_PATH = "/var/lib/libvirt/images"
)

type Libvirt struct {
	conn *libvirt.Connect
}

func NewLibvirt(conn *libvirt.Connect) (*Libvirt, error) {
	return &Libvirt{conn: conn}, nil
}

func NewLibvirtFromConnectUrl(connectUrl string) (*Libvirt, error) {
	conn, err := libvirt.NewConnect(connectUrl)
	if err != nil {
		return nil, err
	}

	return NewLibvirt(conn)
}

func NewLibvirtFromConnectBuilder(connectBuilder *builder.LibvirtConnectBuilder) (*Libvirt, error) {
	conn, err := connectBuilder.Build()
	if err != nil {
		return nil, err
	}

	return NewLibvirt(conn)
}

func (service *Libvirt) GetDomainByName(name string) (*libvirt.Domain, error) {
	return service.conn.LookupDomainByName(name)
}

func (service *Libvirt) DefineDomainFromXMLString(xmlString string) (*libvirt.Domain, error) {
	return service.conn.DomainDefineXML(xmlString)
}

func (service *Libvirt) DefineDomainFromBuilder(domainBuilder *builder.LibvirtDomainBuilder) (*libvirt.Domain, error) {
	return domainBuilder.Build(service.conn)
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
	return service.conn.ListAllDomains(flags)
}

func (service *Libvirt) Cleanup() error {
	_, err := service.conn.Close()
	return err
}
