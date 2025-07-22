package types

import (
	"os/exec"

	"golang.org/x/crypto/ssh"
)

type ShellExecutor interface {
	Exec(string) error
}

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
	session, err := executor.client.NewSession()
	if err != nil {
		return err
	}
	err = session.Run(command)
	return err
}

type LocalShellExecutor struct {
}

func NewLocalShellExecutor() (*LocalShellExecutor, error) {
	return &LocalShellExecutor{}, nil
}

func (executor *LocalShellExecutor) Exec(command string) error {
	cmd := exec.Command(command)
	err := cmd.Run()
	return err
}
