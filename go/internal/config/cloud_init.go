package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/nnurry/harmonia/internal/utils"
)

type userData struct {
	Hostname              string
	SshUser               string
	PublicSshKeysContents []string
	DisableRootPW         bool
}

type metaData struct {
	Hostname   string
	InstanceId string `json:"instance-id"`
}

type networkConfig struct {
	IpV4Addresses    []string
	Nameservers      []string
	GatewayV4Address string
	MacAddress       string
	Dhcp4            bool
}

type cloudInitUser struct {
	Name              string   `yaml:"name"`
	Sudo              string   `yaml:"sudo,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
}

type cloudConfigUserDataYAML struct {
	Hostname       string          `yaml:"hostname"`
	ManageEtcHosts bool            `yaml:"manage_etc_hosts"`
	DisableRootPW  bool            `yaml:"disable_root_pw"`
	Users          []cloudInitUser `yaml:"users"`
}

type networkConfigNameserversYAML struct {
	Addresses []string `yaml:"addresses"`
}

type networkConfigEthernetYAML struct {
	Dhcp4       bool                         `yaml:"dhcp4"`
	Addresses   []string                     `yaml:"addresses,omitempty"`
	Gateway4    string                       `yaml:"gateway4,omitempty"`
	MacAddress  string                       `yaml:"macaddress,omitempty"`
	Nameservers networkConfigNameserversYAML `yaml:"nameservers,omitempty"`
}

type networkConfigNetworkYAML struct {
	Version   int                                  `yaml:"version"`
	Ethernets map[string]networkConfigEthernetYAML `yaml:"ethernets"`
}

type cloudConfigNetworkConfigYAML struct {
	Network networkConfigNetworkYAML `yaml:"network"`
}

type CloudInitConfig struct {
	userData      userData
	metaData      metaData
	networkConfig networkConfig
}

func NewCloudInitConfig(
	hostname string,
	sshUser string,
	sshPublicKeys []string,
	ipv4Address string,
	nameservers []string,
	ipv4GatewayAddress string,
	macAddress string,
	disableRootPW bool,
) *CloudInitConfig {
	ci := &CloudInitConfig{}

	ci.userData = userData{
		Hostname:              hostname,
		SshUser:               sshUser,
		PublicSshKeysContents: sshPublicKeys,
		DisableRootPW:         disableRootPW, // Set this value from input
	}

	instanceID := fmt.Sprintf("monochromatic%s-%s", hostname, uuid.New().String()[:8])
	ci.metaData = metaData{
		Hostname:   hostname,
		InstanceId: instanceID,
	}

	ci.networkConfig = networkConfig{
		IpV4Addresses:    []string{fmt.Sprintf("%s/24", ipv4Address)},
		Nameservers:      nameservers,
		GatewayV4Address: ipv4GatewayAddress,
		MacAddress:       macAddress,
		Dhcp4:            false,
	}

	return ci
}

func (ci *CloudInitConfig) ToUserDataYAML() (string, error) {
	user := cloudInitUser{
		Name:              ci.userData.SshUser,
		Sudo:              "ALL=(ALL) NOPASSWD:ALL",
		SSHAuthorizedKeys: ci.userData.PublicSshKeysContents,
	}

	ciYaml := cloudConfigUserDataYAML{
		Hostname:       ci.userData.Hostname,
		ManageEtcHosts: true,
		DisableRootPW:  ci.userData.DisableRootPW,
		Users:          []cloudInitUser{user},
	}

	return utils.MarshalYAML(&ciYaml, "#cloud-config")
}

func (ci *CloudInitConfig) ToMetaDataJSON() (string, error) {
	byteSlice, err := json.MarshalIndent(ci.metaData, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal meta data to JSON: %w", err)
	}
	return string(byteSlice), nil
}

func (ci *CloudInitConfig) ToNetworkConfigYAML() (string, error) {
	ethernet := networkConfigEthernetYAML{
		Dhcp4:      ci.networkConfig.Dhcp4,
		MacAddress: ci.networkConfig.MacAddress,
	}

	if !ci.networkConfig.Dhcp4 {
		ethernet.Addresses = ci.networkConfig.IpV4Addresses
		ethernet.Gateway4 = ci.networkConfig.GatewayV4Address
	}

	if len(ci.networkConfig.Nameservers) > 0 {
		ethernet.Nameservers = networkConfigNameserversYAML{
			Addresses: ci.networkConfig.Nameservers,
		}
	}

	networkConfigData := cloudConfigNetworkConfigYAML{
		Network: networkConfigNetworkYAML{
			Version: 2,
			Ethernets: map[string]networkConfigEthernetYAML{
				"eth0": ethernet,
			},
		},
	}

	return utils.MarshalYAML(&networkConfigData, "") // No header for network-config
}

func NewCloudInitFromNodeConfig(node VMNode, group SharedConfigs) (*CloudInitConfig, error) {
	publicKeyPath, err := group.SSH.GetExpandedPublicKeyPath()
	if err != nil {
		return nil, fmt.Errorf("error expanding public key path: %w", err)
	}
	publicKeyContent, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file %s: %w", publicKeyPath, err)
	}

	return NewCloudInitConfig(
		node.Name,
		group.SSH.User,
		[]string{string(publicKeyContent)},
		node.IPAddress,
		group.CloudInit.Nameservers,
		node.GatewayAddress,
		node.MACAddress,
		group.CloudInit.DisableRootPW,
	), nil
}
