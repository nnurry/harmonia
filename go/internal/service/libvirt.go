package service

import (
	"libvirt.org/go/libvirt"
)

type LibvirtService struct {
	conn *libvirt.Connect
}

func NewLibvirtService(conn *libvirt.Connect) *LibvirtService {
	return &LibvirtService{conn: conn}
}
