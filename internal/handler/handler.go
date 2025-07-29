package handler

import (
	"net/http"

	"github.com/nnurry/harmonia/internal/contract"
)

func writeResult(writer http.ResponseWriter, code int, result contract.Result) {
	writer.WriteHeader(code)
	writer.Write(result.Compile())
}
