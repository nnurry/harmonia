package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nnurry/harmonia/internal/contract"
	"github.com/rs/zerolog/log"
)

type responseCallback func()

type VirtualMachine struct {
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{}
}

func (handler VirtualMachine) parseBodyAndHandleError(writer http.ResponseWriter, request *http.Request, v any) (responseCallback, error) {
	err := json.NewDecoder(request.Body).Decode(&v)
	if err != nil {
		return func() {
			result := contract.Result{
				Body: struct {
					Error error `json:"error"`
				}{err},
				Message: "could not parse request body",
			}
			writeResult(writer, http.StatusBadRequest, result)
		}, err
	}
	return func() {}, nil
}

func (handler *VirtualMachine) create(request contract.BuildVirtualMachineConfig) error {
	log.Info().Msg(fmt.Sprintf("create VM %v", request.Name))
	// create VM
	// create cloud-init ISO
	return nil
}

func (handler *VirtualMachine) Create(writer http.ResponseWriter, request *http.Request) {
	var createRequest contract.BuildVirtualMachineRequest
	cb, err := handler.parseBodyAndHandleError(writer, request, &createRequest)
	if err != nil {
		cb()
		return
	}

	err = handler.create(contract.BuildVirtualMachineConfig(createRequest))
	if err != nil {
		result := contract.Result{
			Body: struct {
				Name  string `json:"name"`
				Error error  `json:"error"`
			}{
				Name:  createRequest.Name,
				Error: err,
			},
			Message: "could not create single virtual machine",
		}
		writeResult(writer, http.StatusInternalServerError, result)
		return
	}

	writeResult(writer, http.StatusOK, contract.Result{
		Message: "created single virtual machine",
	})
}

func (handler *VirtualMachine) CreateFleet(writer http.ResponseWriter, request *http.Request) {
	var fleetCreateRequest contract.BuildVirtualMachineFleetRequest
	cb, err := handler.parseBodyAndHandleError(writer, request, &fleetCreateRequest)

	if err != nil {
		cb()
		return
	}

	errors := map[string]error{}

	for _, config := range fleetCreateRequest.GetCoalesced().VirtualMachineConfigs {
		err := handler.create(contract.BuildVirtualMachineConfig(config))
		if err != nil {
			errors[config.Name] = err
		}
	}

	if len(errors) > 0 {
		result := contract.Result{
			Body:    errors,
			Message: "could not create virtual machine fleet",
		}
		writeResult(writer, http.StatusInternalServerError, result)
		return
	}
	writeResult(
		writer, http.StatusOK,
		contract.Result{Message: "created virtual machine fleet"},
	)
}
