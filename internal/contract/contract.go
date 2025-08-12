package contract

import (
	"encoding/json"

	"github.com/nnurry/harmonia/internal/logger"
)

type GenericResponse struct {
	Body    any    `json:"body"`
	Message string `json:"message"`
}

func (r GenericResponse) Compile() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		logger.Err(err, "got into error while serializing response")
	}
	return data
}
