package contract

type CloneVirtualMachineRequest struct {
	GeneralVMConfig `yaml:",inline"`
	UserVMConfig    `yaml:",inline"`
	NetworkVMConfig `yaml:",inline"`
}

type GeneralVMConfig struct {
	Name                   string  `yaml:"name"`
	BaseVirtualMachineName string  `yaml:"base_vm_name"`
	NumOfVCPUs             int     `yaml:"vcpu"`
	MemoryInGiB            float64 `yaml:"memory_gb"`
	DiskSizeInGiB          float64 `yaml:"disk_gb"`
	IsCopyOnWriteClone     bool    `yaml:"cow_clone"`
}

type UserVMConfig struct {
	User                  string   `yaml:"user"`
	AuthorizedKeyPaths    []string `yaml:"authorized_key_paths"`
	AuthorizedKeyContents []string `yaml:"authorized_key_contents"`
	DisableRootPassword   bool     `yaml:"disable_root_pw"`
}

type NetworkVMConfig struct {
	IPv4Address        string   `yaml:"ip_address"`
	IPv4GatewayAddress string   `yaml:"gateway_address"`
	MacAddress         string   `yaml:"mac_addresss"`
	Nameservers        []string `yaml:"nameservers"`
}
