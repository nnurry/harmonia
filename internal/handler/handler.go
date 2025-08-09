package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/goccy/go-yaml"
	"github.com/nnurry/harmonia/internal/contract"
)

func writeResult(writer http.ResponseWriter, code int, result contract.GenericResponse) {
	writer.WriteHeader(code)
	writer.Write(result.Compile())
}

func parseBody(request *http.Request, v any, isYAMLable bool) error {
	// prioritize JSON over YAML
	var (
		parseErr  error  = nil
		jsonErr   error  = nil
		yamlErr   error  = nil
		bytesBody []byte = nil
	)

	bytesBody, parseErr = io.ReadAll(request.Body)
	if parseErr != nil {
		return fmt.Errorf("failed to read request body as bytes: %v", parseErr)
	}

	jsonErr = json.NewDecoder(bytes.NewReader(bytesBody)).Decode(v)
	if jsonErr == nil {
		return nil
	} else {
		jsonErr = fmt.Errorf("failed to parse body as JSON: %v", jsonErr)
	}

	if isYAMLable {
		yamlErr = yaml.NewDecoder(bytes.NewReader(bytesBody)).Decode(v)
		if yamlErr == nil {
			return nil
		} else {
			yamlErr = fmt.Errorf("failed to parse body as YAML: %v", yamlErr)
		}
	}

	return errors.Join(jsonErr, yamlErr)
}

func parseBodyAndHandleError(writer http.ResponseWriter, request *http.Request, v any, isYAMLable bool) (responseCallback, error) {
	if err := parseBody(request, v, isYAMLable); err != nil {
		subErrAsStringList := []string{}
		for _, subErr := range err.(interface{ Unwrap() []error }).Unwrap() {
			subErrAsStringList = append(subErrAsStringList, subErr.Error())
		}
		return func() {
			result := contract.GenericResponse{
				Body: struct {
					Errors []string `json:"errors"`
				}{Errors: subErrAsStringList},
				Message: "could not parse request body",
			}
			writeResult(writer, http.StatusBadRequest, result)
		}, err
	}
	return func() {}, nil
}
