package config

type userData struct {
	hostname              string
	sshUser               string
	publicSshKeysContents []string
}

type metaData struct {
	ipAddresses string
	nameservers []string
	gateway     string
	macAddress  string
}

type networkConfig struct {
	hostname   string
	instanceId string
}

type CloudInit struct {
	UserData      userData
	MetaData      metaData
	NetworkConfig networkConfig
}

// hostname: str,
// ssh_user: str,
// ssh_public_keys: list[str],
// ipv4_address: str,
// nameservers: list[str],
// ipv4_gateway_address: str,
// mac_address: str,
func NewCloudInit() CloudInit {
	return NewCloudInit()
}
