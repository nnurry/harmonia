package contract

import "github.com/nnurry/harmonia/internal/connection"

type VirtualMachineConfig struct {
	GeneralVMConfig             `json:",inline"`
	UserVMConfig                `json:",inline"`
	NetworkVMConfig             `json:",inline"`
	*HypervisorConnectionConfig `json:"hypervisor_connection,omitempty"`

	CloudInitISOPath string `json:"cloud_init_iso_path"`
	QCOW2FilePath    string `json:"qcow2_file_path"`
}

type HypervisorConnectionConfig struct {
	connection.LibvirtConfig `json:"libvirt"`
	connection.SSHConfig     `json:"ssh"`
	IsLocalShell             bool `json:"is_local_shell"`
}

type GeneralVMConfig struct {
	Name                   string  `json:"name"`
	BaseVirtualMachineName string  `json:"base_vm_name"`
	NumOfVCPUs             int     `json:"vcpu"`
	MemoryInGiB            float64 `json:"memory_gb"`
	DiskSizeInGiB          float64 `json:"disk_gb"`
	IsCopyOnWriteClone     bool    `json:"is_cow_clone"`
}

type UserVMConfig struct {
	User                  string   `json:"user"`
	AuthorizedKeyPaths    []string `json:"authorized_key_paths"`
	AuthorizedKeyContents []string `json:"authorized_key_contents"`
	DisableRootPassword   bool     `json:"disable_root_pw,omitempty"`
}

type NetworkVMConfig struct {
	IPv4Address        string   `json:"ip_address"`
	IPv4GatewayAddress string   `json:"gateway_address"`
	MacAddress         string   `json:"mac_address"`
	Nameservers        []string `json:"nameservers"`
}

type CreateVirtualMachineRequest struct {
	VirtualMachineConfig `json:",inline"`
}

type CreateVirtualMachineResult struct {
	UUID  string `json:"uuid,omitempty"`
	Name  string `json:"name"`
	Error string `json:"error,omitempty"`
}

type DeleteVirtualMachineResult struct {
	UUID  string `json:"uuid,omitempty"`
	Name  string `json:"name"`
	Error string `json:"error,omitempty"`
}
