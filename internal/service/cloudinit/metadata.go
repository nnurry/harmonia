package cloudinit

import (
	"bytes"
	"encoding/json"

	"github.com/nnurry/harmonia/pkg/utils"
)

type MetaData struct {
	InstanceId string `json:"instance-id,omitempty"`
	Hostname   string `json:"local-hostname,omitempty"`
}

func (md MetaData) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	return utils.SerializeFromEncoder(json.NewEncoder(&buf), &buf, md)
}
