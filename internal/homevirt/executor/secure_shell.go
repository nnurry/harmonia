package executor

import (
	"golang.org/x/crypto/ssh"
)

type SecureShellExecutor struct {
	client  *ssh.Client
	session *ssh.Session
}

func NewSecureShellExecutor(client *ssh.Client) (*SecureShellExecutor, error) {
	executor := &SecureShellExecutor{client: client}
	err := executor.CreateNewSession()
	if err != nil {
		return nil, err
	}
	return executor, nil
}

func (executor *SecureShellExecutor) CreateNewSession() error {
	session, err := executor.client.NewSession()
	if err != nil {
		return err
	}
	executor.session = session
	return nil
}

func (executor *SecureShellExecutor) Exec(command string) error {
	err := executor.session.Run(command)
	return err
}
