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
	writeBytes(writer, code, result.Compile())
}

func writeBytes(writer http.ResponseWriter, code int, bytesData []byte) {
	writer.WriteHeader(code)
	writer.Write(bytesData)
}

func readBodyFromRequestAsBytes(request *http.Request) ([]byte, error) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		err = fmt.Errorf("failed to read request body as bytes: %v", err)
	}
	return body, err
}

func parseJSONFromBytes(bytesBody []byte, v any) error {
	err := json.NewDecoder(bytes.NewReader(bytesBody)).Decode(v)
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to parse body as JSON: %v", err)
}

func parseYAMLFromBytes(bytesBody []byte, v any) error {
	err := yaml.NewDecoder(bytes.NewReader(bytesBody)).Decode(v)
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to parse body as YAML: %v", err)
}

func parseBody(request *http.Request, v any, isYAMLable bool) error {
	// prioritize JSON over YAML
	var (
		parseErr  error  = nil
		jsonErr   error  = nil
		yamlErr   error  = nil
		bytesBody []byte = nil
	)

	bytesBody, parseErr = readBodyFromRequestAsBytes(request)
	if parseErr != nil {
		return parseErr
	}

	jsonErr = parseJSONFromBytes(bytesBody, v)
	if jsonErr == nil {
		return nil
	}

	if isYAMLable {
		yamlErr = parseYAMLFromBytes(bytesBody, v)
		if yamlErr == nil {
			return nil
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
