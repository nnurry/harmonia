package handler

import (
	"encoding/json"
	"net/http"

	"github.com/nnurry/harmonia/internal/connection"
	"github.com/nnurry/harmonia/internal/contract"
	"github.com/nnurry/harmonia/internal/logger"
	"github.com/nnurry/harmonia/internal/processor"
	"github.com/nnurry/harmonia/internal/service"
)

type responseCallback func()

type VirtualMachine struct {
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{}
}

func (handler *VirtualMachine) create(config contract.VirtualMachineConfig) (string, error) {
	var (
		sshConnection    *connection.SSH
		shellProcessor   ShellProcessor
		libvirtService   LibvirtService
		cloudInitService CloudInitService
	)

	if config.HypervisorConnectionConfig.IsLocalShell {
		shellProcessor = processor.NewLocalShell()
	} else {
		var err error
		sshConnection, err = connection.NewSSH(config.HypervisorConnectionConfig.SSHConfig)
		if err != nil {
			return "", err
		}
		shellProcessor = processor.NewSecureShell(sshConnection)
	}

	// create services
	if conn, err := connection.NewLibvirt(config.HypervisorConnectionConfig.LibvirtConfig); err != nil {
		return "", err
	} else {
		libvirtService, err = service.NewLibvirt(conn)
		if err != nil {
			return "", err
		}
	}

	cloudInitService, err := service.NewCloudInit(shellProcessor, sshConnection)
	if err != nil {
		return "", err
	}

	virtualMachineService, _ := service.NewVirtualMachine(libvirtService, cloudInitService, shellProcessor)
	return virtualMachineService.Create(config)
}

func (handler *VirtualMachine) Create(writer http.ResponseWriter, request *http.Request) {
	var createRequest contract.CreateVirtualMachineRequest
	cb, err := parseBodyAndHandleError(writer, request, &createRequest, true)
	if err != nil {
		cb()
		return
	}

	domainUuid, err := handler.create(createRequest.VirtualMachineConfig)
	result := contract.CreateVirtualMachineResult{
		Name: createRequest.Name,
	}
	if err != nil {
		result.Error = err.Error()
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

func (handler *VirtualMachine) FormatRequest(writer http.ResponseWriter, request *http.Request) {
	contractGeneratorMap := map[string]func() any{
		"create":       func() any { return contract.CreateVirtualMachineRequest{} },
		"create_fleet": func() any { return contract.CreateVirtualMachineFleetRequest{} },
	}

	serializerMap := map[string]func(any) ([]byte, error){
		"json": func(v any) ([]byte, error) { return json.MarshalIndent(v, "", " ") },
	}

	queries := request.URL.Query()

	var (
		inputData  any
		outputData []byte
		err        error
	)

	if contractGenerator, ok := contractGeneratorMap[queries.Get("contract")]; !ok {
		writeResult(writer, http.StatusNotFound, contract.GenericResponse{
			Body:    nil,
			Message: "no matching contract to format request",
		})
		return
	} else {
		inputData = contractGenerator()
	}

	cb, err := parseBodyAndHandleError(writer, request, &inputData, true)
	if err != nil {
		cb()
		return
	}

	if serializer, ok := serializerMap[queries.Get("format")]; !ok {
		writeResult(writer, http.StatusNotFound, contract.GenericResponse{
			Body:    nil,
			Message: "no matching serializer to format request",
		})
		return
	} else {
		outputData, err = serializer(inputData)
	}

	if err != nil {
		writeResult(writer, http.StatusInternalServerError, contract.GenericResponse{
			Body:    err,
			Message: "could not serialize data",
		})
		return
	}

	writeBytes(writer, http.StatusOK, outputData)
}

func (handler *VirtualMachine) CreateFleet(writer http.ResponseWriter, request *http.Request) {
	var fleetCreateRequest contract.CreateVirtualMachineFleetRequest
	cb, err := parseBodyAndHandleError(writer, request, &fleetCreateRequest, true)

	if err != nil {
		cb()
		return
	}

	result := contract.CreateVirtualMachineFleetResult{
		SubResults: []contract.CreateVirtualMachineResult{},
		Failed:     0,
		Success:    0,
		Total:      0,
	}

	for _, config := range fleetCreateRequest.GetCoalesced().VirtualMachineConfigs {
		subResult := contract.CreateVirtualMachineResult{
			Name: config.Name,
		}

		logger.Infof("creating VM %v", config.GeneralVMConfig.Name)
		domainUuid, err := handler.create(config)

		if err != nil {
			subResult.Error = err.Error()
			logger.Errorf("failed to create VM %v: %v", config.GeneralVMConfig.Name, subResult.Error)
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
