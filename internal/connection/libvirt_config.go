package connection

type LibvirtConfig struct {
	ConnectionUrl string `json:"connection_url"`
	KeyfilePath   string `json:"keyfile_path"`
}
