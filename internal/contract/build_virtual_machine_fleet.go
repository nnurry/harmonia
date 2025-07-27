package contract

type BuildVirtualMachineFleetRequest struct {
	SharedConfig          FleetSharedConfig            `yaml:"shared_config"`
	VirtualMachineConfigs []CloneVirtualMachineRequest `yaml:"virtual_machines"`
}

type FleetSharedConfig struct {
	SSH     SSHSharedConfig     `yaml:"ssh"`
	Network NetworkSharedConfig `yaml:"cloud_init"`
}

type SSHSharedConfig struct {
	User                  string   `yaml:"user"`
	DisableRootPassword   bool     `yaml:"disable_root_pw"`
	AuthorizedKeyPaths    []string `yaml:"authorized_key_paths"`
	AuthorizedKeyContents []string `yaml:"authorized_key_contents"`
}

type NetworkSharedConfig struct {
	Nameservers []string `yaml:"nameservers"`
}
