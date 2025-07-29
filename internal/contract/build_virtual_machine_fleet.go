package contract

type BuildVirtualMachineFleetRequest struct {
	SharedConfig          FleetSharedConfig           `json:"shared_config"`
	VirtualMachineConfigs []BuildVirtualMachineConfig `json:"virtual_machines"`
}

type FleetSharedConfig struct {
	GeneralSharedConfig `json:"general"`
	SSHSharedConfig     `json:"ssh"`
	NetworkSharedConfig `json:"cloud_init"`
}

type GeneralSharedConfig struct {
	BaseVirtualMachineName string `json:"base_vm_name"`
}

type SSHSharedConfig struct {
	User                  string   `json:"user"`
	AuthorizedKeyPaths    []string `json:"authorized_key_paths"`
	AuthorizedKeyContents []string `json:"authorized_key_contents"`
}

type NetworkSharedConfig struct {
	Nameservers []string `json:"nameservers"`
}

func (r BuildVirtualMachineFleetRequest) GetCoalesced() BuildVirtualMachineFleetRequest {
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
	}

	return r
}

type BuildVirtualMachineFleetResult struct {
	SubResults []BuildVirtualMachineResult `json:"sub_results"`
	Failed     int                         `json:"failed"`
	Success    int                         `json:"success"`
	Total      int                         `json:"total"`
}
