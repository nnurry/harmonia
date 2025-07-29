package handler

import (
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

func (handler *VirtualMachine) create(request contract.BuildVirtualMachineConfig) (string, error) {
	log.Info().Msg(fmt.Sprintf("create VM %v", request.Name))
	// create VM
	// create cloud-init ISO
	return "", nil
}

func (handler *VirtualMachine) Create(writer http.ResponseWriter, request *http.Request) {
	var createRequest contract.BuildVirtualMachineRequest
	cb, err := parseBodyAndHandleError(writer, request, &createRequest)
	if err != nil {
		cb()
		return
	}

	domainUuid, err := handler.create(contract.BuildVirtualMachineConfig(createRequest))
	result := contract.BuildVirtualMachineResult{
		Name: createRequest.Name,
	}
	if err != nil {
		result.Error = err
		writeResult(writer, http.StatusInternalServerError, contract.GenericResponse{
			Body:    result,
			Message: "could not create single virtual machine",
		})
		return
	}

	result.UUID = domainUuid

	writeResult(writer, http.StatusOK, contract.GenericResponse{
		Body:    result,
		Message: "created single virtual machine",
	})
}

func (handler *VirtualMachine) CreateFleet(writer http.ResponseWriter, request *http.Request) {
	var fleetCreateRequest contract.BuildVirtualMachineFleetRequest
	cb, err := parseBodyAndHandleError(writer, request, &fleetCreateRequest)

	if err != nil {
		cb()
		return
	}

	result := contract.BuildVirtualMachineFleetResult{
		SubResults: []contract.BuildVirtualMachineResult{},
		Failed:     0,
		Success:    0,
		Total:      0,
	}

	for _, config := range fleetCreateRequest.GetCoalesced().VirtualMachineConfigs {
		subResult := contract.BuildVirtualMachineResult{
			Name: config.Name,
		}

		domainUuid, err := handler.create(contract.BuildVirtualMachineConfig(config))

		if err != nil {
			subResult.Error = err
			result.Failed++
		} else {
			subResult.UUID = domainUuid
			result.Success++
		}
		result.Total++

		result.SubResults = append(result.SubResults, subResult)
	}

	var message string
	if result.Failed > 0 {
		if result.Failed == result.Total {
			message = "failed to create virtual machine fleet"
		} else {
			message = "created virtual machine fleet with partial failures"
		}
	} else {
		message = "created virtual machine fleet"
	}

	writeResult(
		writer, http.StatusOK,
		contract.GenericResponse{
			Body:    result,
			Message: message,
		},
	)
}
