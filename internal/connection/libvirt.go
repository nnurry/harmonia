package connection

import (
	"fmt"
	"net/url"

	"libvirt.org/go/libvirt"
)

type Libvirt struct {
	connect *libvirt.Connect
	url     string
}

func NewLibvirt(config LibvirtConfig) (*Libvirt, error) {
	connection := &Libvirt{}
	if config.ConnectionUrl == "" {
		return nil, fmt.Errorf("empty connection url")
	}

	url, err := url.Parse(config.ConnectionUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid connection url: %v", err)
	}

	query := url.Query()

	if config.KeyfilePath != "" {
		query.Add("keyfile", config.KeyfilePath)
	}

	url.RawQuery = query.Encode()
	connection.url = url.String()

	libvirtConnection, err := libvirt.NewConnect(connection.url)
	if err != nil {
		return nil, fmt.Errorf("could not create libvirt connection: %v", err)
	}

	connection.connect = libvirtConnection

	return connection, nil
}

func (connection Libvirt) URL() string {
	return connection.url
}

func (connection Libvirt) Name() string {
	return "ssh"
}

func (connection *Libvirt) Connect() *libvirt.Connect {
	return connection.connect
}

func (connection *Libvirt) Cleanup() error {
	_, err := connection.connect.Close()
	return err
}
