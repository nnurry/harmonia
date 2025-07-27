package cloudinit

import (
	"bytes"
	"encoding/json"

	"github.com/nnurry/harmonia/pkg/utils"
)

type MetaData struct {
	InstanceId string `yaml:"instance-id,omitempty"`
	Hostname   string `yaml:"local-hostname,omitempty"`
}

func (md MetaData) FileName() string {
	return "meta-data"
}

func (md MetaData) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	return utils.SerializeFromEncoder(json.NewEncoder(&buf), &buf, md)
}
