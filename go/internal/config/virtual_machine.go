package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/nnurry/harmonia/internal/utils"
)

type Config struct {
	Groups []VMGroup `yaml:"groups"`
}

type VMGroup struct {
	Name          string        `yaml:"name"`
	BaseVMName    string        `yaml:"base_vm_name"`
	SharedConfigs SharedConfigs `yaml:"shared_configs"`
	Nodes         []VMNode      `yaml:"nodes"`
}

type SSHConfig struct {
	User           string `yaml:"user"`
	PublicKeyPath  string `yaml:"public_key_path"`
	PrivateKeyPath string `yaml:"private_key_path"`
}

type CloudInitBase struct {
	Nameservers   []string `yaml:"nameservers"`
	DisableRootPW bool     `yaml:"disable_root_pw"`
}

type VMNode struct {
	Name           string `yaml:"name"`
	IPAddress      string `yaml:"ip_address"`
	GatewayAddress string `yaml:"gateway_address"`
	VCPU           int    `yaml:"vcpu"`
	MemoryGB       int    `yaml:"memory_gb"`
	DiskGB         int    `yaml:"disk_gb"`
	MACAddress     string `yaml:"mac_address"`
	IsCowClone     bool   `yaml:"is_cow_clone"`
}

type SharedConfigs struct {
	SSH       SSHConfig     `yaml:"ssh"`
	CloudInit CloudInitBase `yaml:"cloud_init"`
}

func (s *SSHConfig) GetExpandedPublicKeyPath() (string, error) {
	return utils.ExpandPath(s.PublicKeyPath)
}

func LoadConfigFromFile(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML config from %s: %w", filePath, err)
	}

	return &cfg, nil
}
