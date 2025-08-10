package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nnurry/harmonia/internal/connection"
	"github.com/nnurry/harmonia/internal/service/cloudinit"
	"github.com/pkg/sftp"
	"github.com/rs/zerolog/log"
)

const (
	DEFAULT_CLOUD_INIT_ISO_BASE_PATH = "/var/lib/libvirt/images"
)

type CloudInit struct {
	processor     ShellProcessor
	sftpClient    *sftp.Client
	UserData      cloudinit.UserData
	MetaData      cloudinit.MetaData
	NetworkConfig cloudinit.NetworkConfig
}

type cloudInitISOIngredient interface {
	Serialize() ([]byte, error)
	FileName() string
}

func NewCloudInit(processor ShellProcessor, sshConnection *connection.SSH) (*CloudInit, error) {
	log.Info().Msg("creating SFTP client")

	sftpClient, err := sftp.NewClient(sshConnection.Client())
	if err != nil {
		return nil, fmt.Errorf("could not create SFTP client for cloud-init service: %v", err)
	}

	log.Info().Msg("created SFTP client")
	return &CloudInit{processor: processor, sftpClient: sftpClient}, nil
}

func (service *CloudInit) SetUserData(userData cloudinit.UserData) {
	service.UserData = userData
}

func (service *CloudInit) SetMetaData(metaData cloudinit.MetaData) {
	service.MetaData = metaData
}

func (service *CloudInit) SetNetworkConfig(networkConfig cloudinit.NetworkConfig) {
	service.NetworkConfig = networkConfig
}

func (service *CloudInit) WriteToDisk(ctx context.Context, basePath string, filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("empty file path for cloud-init ISO")
	}

	var err error

	if service.sftpClient != nil {
		err = service.sftpClient.MkdirAll(basePath)
		err = errors.Join(err, service.sftpClient.Chmod(basePath, os.FileMode(0777)))
	} else {
		err = os.MkdirAll(basePath, os.FileMode(0777))
	}

	if err != nil {
		return "", fmt.Errorf("could not mkdir '%v': %v", basePath, err)
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

		if service.sftpClient != nil {
			var sftpFile *sftp.File
			if sftpFile, err = service.sftpClient.Create(path); err == nil {
				sftpFile.Chmod(os.FileMode(0777))
				_, err = sftpFile.Write(data)
			}
		} else {
			err = os.WriteFile(path, data, os.FileMode(0777))
		}

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

	err = service.processor.Execute(ctx, os.Stdout, os.Stderr, cmdParts[0], cmdParts[1:]...)
	if err != nil {
		return "", fmt.Errorf("could not write ISO file to disk for cloud-init ISO: %v", err)
	}
	return isoFilePath, nil
}

func (service *CloudInit) RemoveFromDisk(ctx context.Context, basePath string) error {
	return service.sftpClient.RemoveAll(basePath)
}
