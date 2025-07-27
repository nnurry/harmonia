package service

import (
	"context"
	"fmt"
	"os"

	"github.com/nnurry/harmonia/internal/processor"
	"github.com/nnurry/harmonia/internal/service/cloudinit"
)

const (
	DEFAULT_CLOUD_INIT_ISO_BASE_PATH = "/var/lib/libvirt/images"
)

type CloudInit struct {
	processor     processor.Shell
	UserData      cloudinit.UserData
	MetaData      cloudinit.MetaData
	NetworkConfig cloudinit.NetworkConfig
}

type cloudInitISOIngredient interface {
	Serialize() ([]byte, error)
	FileName() string
}

func NewCloudInit() (*CloudInit, error) {
	return &CloudInit{}, nil
}

func (service *CloudInit) SerializeUserData() ([]byte, error) {
	return service.UserData.Serialize()
}

func (service *CloudInit) SerializeNetworkConfig() ([]byte, error) {
	return service.NetworkConfig.Serialize()
}

func (service *CloudInit) SerializeMetadata() ([]byte, error) {
	return service.MetaData.Serialize()
}

func (service *CloudInit) WriteToDisk(ctx context.Context, basePath string, filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("empty file path for cloud-init ISO")
	}

	isoFilePath := fmt.Sprintf("%v/%v", basePath, filename)
	paths := make([]string, 3)

	for _, ingredient := range []cloudInitISOIngredient{
		service.UserData,
		service.MetaData,
		service.NetworkConfig,
	} {
		path := fmt.Sprintf("%v/%v", basePath, ingredient.FileName())
		name := ingredient.FileName()

		data, err := ingredient.Serialize()
		if err != nil {
			return "", fmt.Errorf("could not serialize ingredient %v for cloud-init ISO: %v", name, err)
		}

		err = os.WriteFile(path, data, os.FileMode(0777))
		if err != nil {
			return "", fmt.Errorf("could not write ingredient %v to disk for cloud-init ISO: %v", name, err)
		}

		paths = append(paths, path)
	}

	cmdParts := []string{
		"mkisofs",
		"-output", isoFilePath,
		"-volid", "cidata",
		"-joliet",
		"-r",
	}
	cmdParts = append(cmdParts, paths...)

	err := service.processor.Execute(ctx, os.Stdout, os.Stderr, cmdParts[0], cmdParts[1:]...)
	if err != nil {
		return "", fmt.Errorf("could not write ISO file to disk for cloud-init ISO: %v", err)
	}
	return isoFilePath, nil
}
