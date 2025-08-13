package contract

import "fmt"

type VirtualMachineFleetConfig struct {
	SharedConfig          FleetSharedConfig      `json:"shared_config"`
	VirtualMachineConfigs []VirtualMachineConfig `json:"virtual_machines"`
}

type FleetSharedConfig struct {
	GeneralSharedConfig         `json:"general"`
	SSHSharedConfig             `json:"ssh"`
	NetworkSharedConfig         `json:"cloud_init"`
	*HypervisorConnectionConfig `json:"hypervisor_connection,omitempty"`
}

type GeneralSharedConfig struct {
	BaseVirtualMachineName  string `json:"base_vm_name"`
	VirtualMachineFleetName string `json:"fleet_name"`
}

type SSHSharedConfig struct {
	User                  string   `json:"user"`
	AuthorizedKeyPaths    []string `json:"authorized_key_paths"`
	AuthorizedKeyContents []string `json:"authorized_key_contents"`
}

type NetworkSharedConfig struct {
	Nameservers []string `json:"nameservers"`
}

func (r VirtualMachineFleetConfig) GetCoalesced() VirtualMachineFleetConfig {
	for i, vmConfig := range r.VirtualMachineConfigs {
		if len(vmConfig.Nameservers) < 1 {
			r.VirtualMachineConfigs[i].Nameservers = r.SharedConfig.Nameservers
		}

		if len(vmConfig.AuthorizedKeyPaths) < 1 {
			r.VirtualMachineConfigs[i].AuthorizedKeyPaths = r.SharedConfig.AuthorizedKeyPaths
		}

		if len(vmConfig.AuthorizedKeyContents) < 1 {
			r.VirtualMachineConfigs[i].AuthorizedKeyContents = r.SharedConfig.AuthorizedKeyContents
		}

		if vmConfig.User == "" {
			r.VirtualMachineConfigs[i].User = r.SharedConfig.User
		}

		if vmConfig.BaseVirtualMachineName == "" {
			r.VirtualMachineConfigs[i].BaseVirtualMachineName = r.SharedConfig.BaseVirtualMachineName
		}

		if vmConfig.HypervisorConnectionConfig == nil {
			sharedHypervisorConnectionConfig := *r.SharedConfig.HypervisorConnectionConfig
			r.VirtualMachineConfigs[i].HypervisorConnectionConfig = &sharedHypervisorConnectionConfig
		}

		if r.SharedConfig.GeneralSharedConfig.VirtualMachineFleetName != "" {
			r.VirtualMachineConfigs[i].GeneralVMConfig.Name = fmt.Sprintf(
				"%v-%v",
				r.SharedConfig.GeneralSharedConfig.VirtualMachineFleetName,
				r.VirtualMachineConfigs[i].GeneralVMConfig.Name,
			)
		}
	}

	return r
}

type CreateVirtualMachineFleetRequest struct {
	VirtualMachineFleetConfig `json:",inline"`
}

type CreateVirtualMachineFleetResult struct {
	SubResults []CreateVirtualMachineResult `json:"sub_results"`
	Failed     int                          `json:"failed"`
	Success    int                          `json:"success"`
	Total      int                          `json:"total"`
}

// Same as CreateVirtualMachineFleetRequest
type DeleteVirtualMachineFleetRequest struct {
	VirtualMachineFleetConfig `json:",inline"`
}

type DeleteVirtualMachineFleetResult struct {
	SubResults []DeleteVirtualMachineResult `json:"sub_results"`
	Failed     int                          `json:"failed"`
	Success    int                          `json:"success"`
	Total      int                          `json:"total"`
}
