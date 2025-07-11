package cloudinit

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nnurry/harmonia/internal/config"
	"github.com/nnurry/harmonia/internal/utils"
)

type Generator struct {
	metaData      metaData
	userData      userDataYAML
	networkConfig networkConfigYAML
}

type metaData struct {
	InstanceId string `json:"instance-id"`
	Hostname   string `json:"hostname"`
}

type userDataYAML struct {
	Hostname       string `yaml:"hostname"`
	ManageEtcHosts bool   `yaml:"manage_etc_hosts"`
	DisableRootPW  bool   `yaml:"disable_root_pw"`
	Users          []user `yaml:"users"`
}

type user struct {
	Name              string   `yaml:"name"`
	Sudo              string   `yaml:"sudo"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
}

type networkConfigYAML struct {
	Network struct {
		Version   int                 `yaml:"version"`
		Ethernets map[string]ethernet `yaml:"ethernets"`
	} `yaml:"network"`
}

type ethernet struct {
	Dhcp4       bool     `yaml:"dhcp4"`
	MacAddress  string   `yaml:"macaddress,omitempty"`
	Addresses   []string `yaml:"addresses,omitempty"`
	Gateway4    string   `yaml:"gateway4,omitempty"`
	Nameservers struct {
		Addresses []string `yaml:"addresses,omitempty"`
	} `yaml:"nameservers"`
}

func NewGenerator(node config.VMNode, shared config.SharedConfigs, sshKeyContent string) *Generator {
	instanceID := fmt.Sprintf("monochromatic-%s-%s", node.Name, uuid.NewString()[:8])

	md := metaData{
		InstanceId: instanceID,
		Hostname:   node.Name,
	}

	ud := userDataYAML{
		Hostname:       node.Name,
		ManageEtcHosts: true,
		DisableRootPW:  shared.CloudInit.DisableRootPW,
		Users: []user{
			{
				Name:              shared.SSH.User,
				Sudo:              "ALL=(ALL) NOPASSWD:ALL",
				SSHAuthorizedKeys: []string{sshKeyContent},
			},
		},
	}

	eth := ethernet{
		Dhcp4:      false,
		MacAddress: node.MACAddress,
		Addresses:  []string{fmt.Sprintf("%s/24", node.IPAddress)},
		Gateway4:   node.GatewayAddress,
	}
	eth.Nameservers.Addresses = shared.CloudInit.Nameservers

	nc := networkConfigYAML{}
	nc.Network.Version = 2
	nc.Network.Ethernets = map[string]ethernet{"eth0": eth}

	return &Generator{
		metaData:      md,
		userData:      ud,
		networkConfig: nc,
	}
}

func (g *Generator) MetaData() (string, error) {
	bytes, err := json.MarshalIndent(g.metaData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal meta-data: %w", err)
	}
	return string(bytes), nil
}

func (g *Generator) UserData() (string, error) {
	return utils.MarshalYAML(g.userData, "#cloud-config")
}

func (g *Generator) NetworkConfig() (string, error) {
	return utils.MarshalYAML(g.networkConfig, "")
}
