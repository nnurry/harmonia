package service

import (
	"github.com/nnurry/harmonia/internal/homevirt/builder"
	"libvirt.org/go/libvirt"
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

func (service *Libvirt) Cleanup() error {
	_, err := service.conn.Close()
	return err
}
