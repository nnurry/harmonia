package contract

type BuildVirtualMachineFleet struct {
	SharedConfig          FleetSharedConfig            `json:"shared_config"`
	VirtualMachineConfigs []CloneVirtualMachineRequest `json:"virtual_machines"`
}

type FleetSharedConfig struct {
	SSH     SSHSharedConfig     `json:"ssh"`
	Network NetworkSharedConfig `json:"cloud_init"`
}

type SSHSharedConfig struct {
	User                  string   `json:"user"`
	DisableRootPassword   bool     `json:"disable_root_pw"`
	AuthorizedKeyPaths    []string `json:"authorized_key_paths"`
	AuthorizedKeyContents []string `json:"authorized_key_contents"`
}

type NetworkSharedConfig struct {
	Nameservers []string `json:"nameservers"`
}
