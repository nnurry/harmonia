package cloudinit

import (
	"bytes"

	"github.com/goccy/go-yaml"
	"github.com/nnurry/harmonia/pkg/utils"
)

type NetworkConfig struct {
	Network network `yaml:"network"`
}

type network struct {
	Version   int      `yaml:"version"`
	Ethernets ethernet `yaml:"ethernets"`
}

type ethernet struct {
	Eth0 eth0 `yaml:"eth0"`
}

type eth0 struct {
	Dhcp4              bool         `yaml:"dhcp4"`
	IPv4Addresses      []string     `yaml:"addresses"`
	IPv4GatewayAddress string       `yaml:"gateway4"`
	Nameservers        []nameserver `yaml:"nameservers"`
}

type nameserver struct {
	Addresses []string `yaml:"addresses"`
}

func (nc NetworkConfig) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	return utils.SerializeFromEncoder(yaml.NewEncoder(&buf, yaml.Flow(false)), &buf, nc)
}
