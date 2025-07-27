package cloudinit

import (
	"bytes"

	"github.com/goccy/go-yaml"
	"github.com/nnurry/harmonia/pkg/utils"
)

type UserData struct {
	Hostname       string `yaml:"hostname"`
	ManageEtcHosts bool   `yaml:"manage_etc_hosts,omitempty"`
	DisableRootPw  bool   `yaml:"disable_root_pw,omitempty"`
	Users          []User `yaml:"users,omitempty"`
}

type User struct {
	Name           string   `yaml:"name"`
	Sudo           string   `yaml:"sudo,omitempty"`
	AuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
}

func (ud UserData) FileName() string {
	return "user-data"
}

func (ud UserData) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	data, err := utils.SerializeFromEncoder(yaml.NewEncoder(&buf, yaml.Flow(false)), &buf, ud)
	if err != nil {
		return nil, err
	}

	buf = *bytes.NewBuffer([]byte("#cloud-config\n"))
	if _, err = buf.Write(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
