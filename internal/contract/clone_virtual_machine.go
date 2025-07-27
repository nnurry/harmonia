package contract

type CloneVirtualMachineRequest struct {
	GeneralVMConfig
	UserVMConfig
	NetworkVMConfig
}

type GeneralVMConfig struct {
	Name                   string  `json:"name"`
	BaseVirtualMachineName string  `json:"base_vm_name"`
	NumOfVCPUs             int     `json:"vcpu"`
	MemoryInGiB            float64 `json:"memory_gb"`
	DiskSizeInGiB          float64 `json:"disk_db"`
	IsCopyOnWriteClone     bool    `json:"cow_clone"`
}

type UserVMConfig struct {
	User                  string   `json:"user"`
	AuthorizedKeyPaths    []string `json:"authorized_key_paths"`
	AuthorizedKeyContents []string `json:"authorized_key_contents"`
	DisableRootPassword   bool     `json:"disable_root_pw"`
}

type NetworkVMConfig struct {
	IPv4Address        string   `json:"ip_address"`
	IPv4GatewayAddress string   `json:"gateway_address"`
	MacAddress         string   `json:"mac_addresss"`
	Nameservers        []string `json:"nameservers"`
}
