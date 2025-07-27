package cloudinit

import (
	"bytes"

	"github.com/goccy/go-yaml"
	"github.com/nnurry/harmonia/pkg/utils"
)

type NetworkConfig struct {
	Network Network `json:"network"`
}

type Network struct {
	Version   int      `json:"version"`
	Ethernets Ethernet `json:"ethernets"`
}

type Ethernet struct {
	Eth0 Eth0 `json:"eth0"`
}

type Eth0 struct {
	Dhcp4              bool         `json:"dhcp4"`
	IPv4Addresses      []string     `json:"addresses"`
	IPv4GatewayAddress string       `json:"gateway4"`
	Nameservers        []Nameserver `json:"nameservers"`
}

type Nameserver struct {
	Addresses []string `json:"addresses"`
}

func (nc NetworkConfig) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	return utils.SerializeFromEncoder(yaml.NewEncoder(&buf, yaml.Flow(false)), &buf, nc)
}
