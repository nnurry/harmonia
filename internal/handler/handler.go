package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nnurry/harmonia/internal/contract"
)

func writeResult(writer http.ResponseWriter, code int, result contract.GenericResponse) {
	writer.WriteHeader(code)
	writer.Write(result.Compile())
}

func parseBodyAndHandleError(writer http.ResponseWriter, request *http.Request, v any) (responseCallback, error) {
	err := json.NewDecoder(request.Body).Decode(&v)
	if err != nil {
		return func() {
			body := struct {
				Error error `json:"error"`
			}{Error: err}
			result := contract.GenericResponse{
				Body:    body,
				Message: "could not parse request body",
			}
			writeResult(writer, http.StatusBadRequest, result)
		}, err
	}
	return func() {}, nil
}
