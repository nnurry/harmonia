package contract

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type GenericResponse struct {
	Body    any    `json:"body"`
	Message string `json:"message"`
}

func (r GenericResponse) Compile() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		log.Err(err).Msg("got into error while serializing response")
	}
	return data
}
