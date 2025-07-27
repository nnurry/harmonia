package service

import (
	"github.com/nnurry/harmonia/internal/service/cloudinit"
)

const (
	DEFAULT_CLOUD_INIT_ISO_BASE_PATH = "/var/lib/libvirt/images"
)

type CloudInit struct {
	UserData      cloudinit.UserData
	MetaData      cloudinit.MetaData
	NetworkConfig cloudinit.NetworkConfig
}

func NewCloudInit() (*CloudInit, error) {
	return &CloudInit{}, nil
}

func (service *CloudInit) SerializeUserData() ([]byte, error) {
	return service.UserData.Serialize()
}

func (service *CloudInit) SerializeNetworkConfig() ([]byte, error) {
	return service.NetworkConfig.Serialize()
}

func (service *CloudInit) SerializeMetadata() ([]byte, error) {
	return service.MetaData.Serialize()
}
