package service

const (
	DEFAULT_CLOUD_INIT_ISO_BASE_PATH = "/var/lib/libvirt/images"
)

type CloudInit struct {
}

func NewCloudInit() (*CloudInit, error) {
	return &CloudInit{}, nil
}
