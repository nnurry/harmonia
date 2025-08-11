package cloudinit

import (
	"bytes"

	"github.com/goccy/go-yaml"
	"github.com/nnurry/harmonia/pkg/utils"
)

type NetworkConfig struct {
	Network Network `yaml:"network"`
}

type Network struct {
	Version   int      `yaml:"version"`
	Ethernets Ethernet `yaml:"ethernets"`
}

type Ethernet struct {
	Eth0 Eth0 `yaml:"eth0"`
}

type Eth0 struct {
	Dhcp4              bool         `yaml:"dhcp4"`
	IPv4Addresses      []string     `yaml:"addresses,flow"`
	IPv4GatewayAddress string       `yaml:"gateway4"`
	Nameservers        []Nameserver `yaml:"nameservers"`
}

type Nameserver struct {
	Addresses []string `yaml:"addresses"`
}

func (nc NetworkConfig) FileName() string {
	return "network-config"
}

func (nc NetworkConfig) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	return utils.SerializeFromEncoder(yaml.NewEncoder(&buf, yaml.Flow(false)), &buf, nc)
}
