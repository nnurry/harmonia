package service

import (
	"bytes"
	"fmt"
	"os"

	"github.com/nnurry/harmonia/internal/config"
	"golang.org/x/crypto/ssh"
)

type SSH struct {
	client *ssh.Client
}

func NewSSH(config config.SSH) (*SSH, error) {
	service := &SSH{}
	clientConfig := ssh.ClientConfig{
		User:            config.User,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: config.HostKeyCallback,
	}

	var authMethod ssh.AuthMethod
	var privkeyAuthParseErr, passwordAuthParseErr error

	authMethod, privkeyAuthParseErr = service.ParsePrivateKeyAuth(config)
	if authMethod == nil {
		authMethod, passwordAuthParseErr = service.ParsePassword(config)
	}

	if privkeyAuthParseErr != nil && passwordAuthParseErr != nil {
		return nil, fmt.Errorf("privkey='%v' + password='%v'", privkeyAuthParseErr, passwordAuthParseErr)
	}

	clientConfig.Auth = append(clientConfig.Auth, authMethod)

	dialNetwork := "tcp"
	dialAddress := bytes.NewBufferString(config.Host)
	if config.Port > 0 {
		fmt.Fprintf(dialAddress, ":%d", config.Port)
	}

	client, err := ssh.Dial(dialNetwork, dialAddress.String(), &clientConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create ssh client: %v", err)
	}

	service.client = client

	return service, nil
}

func (service *SSH) ParsePassword(config config.SSH) (ssh.AuthMethod, error) {
	if config.PasswordAuth.Password == "" {
		return nil, fmt.Errorf("empty password")
	}
	return ssh.Password(config.PasswordAuth.Password), nil
}

func (service *SSH) ParsePrivateKeyAuth(config config.SSH) (ssh.AuthMethod, error) {
	if config.PrivateKeyAuth.PrivateKeyPath == "" {
		return nil, fmt.Errorf("empty private key path")
	}

	keyContent, err := os.ReadFile(config.PrivateKeyAuth.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("can't read private key: %v", err)
	}

	var signer ssh.Signer

	if config.PrivateKeyAuth.Passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyContent, []byte(config.PrivateKeyAuth.Passphrase))
		if err != nil {
			return nil, fmt.Errorf("can't parse private key with passphrase: %v", err)
		}
	} else {
		signer, err = ssh.ParsePrivateKey(keyContent)
		if err != nil {
			return nil, fmt.Errorf("can't parse private key: %v", err)
		}
	}

	return ssh.PublicKeys(signer), nil
}

func (service *SSH) Client() *ssh.Client {
	return service.client
}

func (service *SSH) Cleanup() error {
	return service.client.Close()
}
