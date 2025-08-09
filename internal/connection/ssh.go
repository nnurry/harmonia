package connection

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	client *ssh.Client
}

func NewSSH(config SSHConfig) (*SSH, error) {
	connection := &SSH{}
	hostKeyCallback, err := config.HostKeyCallback(config.HostKeyCallbackName)

	if err != nil {
		return nil, fmt.Errorf("could not parse ssh config: %v", err)
	}

	clientConfig := ssh.ClientConfig{
		User:            config.User,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: hostKeyCallback,
		Timeout:         time.Duration(360 * time.Second),
	}

	// let's parse both and see what we got
	privkeyAuthMethod, privkeyAuthParseErr := config.ParsePrivateKeyAuth()
	passwordAuthMethod, passwordAuthParseErr := config.ParsePasswordAuth()

	if privkeyAuthParseErr != nil && passwordAuthParseErr != nil {
		return nil, fmt.Errorf("privkey='%v' + password='%v'", privkeyAuthParseErr, passwordAuthParseErr)
	}

	// prioritize privkey auth over password auth
	if privkeyAuthMethod != nil {
		clientConfig.Auth = append(clientConfig.Auth, privkeyAuthMethod)
	} else {
		clientConfig.Auth = append(clientConfig.Auth, passwordAuthMethod)
	}

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

func (connection SSH) Name() string {
	return "ssh"
}

func (connection *SSH) Client() *ssh.Client {
	return connection.client
}

func (connection *SSH) Session() (*ssh.Session, error) {
	return connection.Client().NewSession()
}

func (connection *SSH) Cleanup() error {
	return connection.client.Close()
}
