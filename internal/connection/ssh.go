package connection

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	client *ssh.Client
}

func NewSSH(config SSHConfig) (*SSH, error) {
	connection := &SSH{}
	clientConfig := ssh.ClientConfig{
		User:            config.User,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: config.HostKeyCallback,
	}

	var authMethod ssh.AuthMethod
	var privkeyAuthParseErr, passwordAuthParseErr error

	authMethod, privkeyAuthParseErr = config.ParsePrivateKeyAuth()
	if authMethod == nil {
		authMethod, passwordAuthParseErr = config.ParsePassword()
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

	connection.client = client

	return connection, nil
}

func (connection *SSH) Client() *ssh.Client {
	return connection.client
}

func (connection *SSH) Cleanup() error {
	return connection.client.Close()
}
