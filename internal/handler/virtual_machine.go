package handler

import (
	"net/http"

	"github.com/nnurry/harmonia/internal/connection"
	"github.com/nnurry/harmonia/internal/contract"
	"github.com/nnurry/harmonia/internal/processor"
	"github.com/nnurry/harmonia/internal/service"
	"github.com/rs/zerolog/log"
)

type responseCallback func()

type VirtualMachine struct {
}

func NewVirtualMachine() *VirtualMachine {
	return &VirtualMachine{}
}

func (handler *VirtualMachine) create(config contract.BuildVirtualMachineConfig) (string, error) {
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
	var createRequest contract.BuildVirtualMachineRequest
	cb, err := parseBodyAndHandleError(writer, request, &createRequest, true)
	if err != nil {
		cb()
		return
	}

	domainUuid, err := handler.create(contract.BuildVirtualMachineConfig(createRequest))
	result := contract.BuildVirtualMachineResult{
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

func (handler *VirtualMachine) CreateFleet(writer http.ResponseWriter, request *http.Request) {
	var fleetCreateRequest contract.BuildVirtualMachineFleetRequest
	cb, err := parseBodyAndHandleError(writer, request, &fleetCreateRequest, true)

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

		log.Info().Msgf("creating VM %v", config.GeneralVMConfig.Name)
		domainUuid, err := handler.create(config)

		if err != nil {
			subResult.Error = err.Error()
			log.Error().Msgf("failed to create VM %v: %v", config.GeneralVMConfig.Name, subResult.Error)
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
