package service

import (
	"libvirt.org/go/libvirt"
)

type Libvirt struct {
	conn *libvirt.Connect
}

func NewLibvirt(conn *libvirt.Connect) *Libvirt {
	return &Libvirt{conn: conn}
}
