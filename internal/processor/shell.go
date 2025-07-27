package processor

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Shell interface {
	Name() string
	Execute(context.Context, io.Writer, io.Writer, string, ...string) error
}

type LocalShell struct {
}

func NewLocalShell() *LocalShell {
	return &LocalShell{}
}

func (processor *LocalShell) Name() string {
	return "local-shell"
}

func (processor *LocalShell) Execute(
	ctx context.Context,
	stdout, stderr io.Writer,
	command string, args ...string,
) error {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("cmd.Run() error: %v", err)
	}
	return nil
}

type SecureShell struct {
	client *ssh.Client
}

func NewSecureShell(client *ssh.Client) *SecureShell {
	return &SecureShell{client: client}
}

func (processor *SecureShell) Name() string {
	return "secure-shell"
}

func (processor *SecureShell) Execute(
	ctx context.Context,
	stdout, stderr io.Writer,
	command string, args ...string,
) error {
	session, err := processor.client.NewSession()
	if err != nil {
		return fmt.Errorf("client.NewSession() error: %v", err)
	}

	defer session.Close()

	commandParts := []string{command}
	commandParts = append(commandParts, args...)

	command = strings.Join(commandParts, " ")

	session.Stdout = stdout
	session.Stderr = stderr

	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("session.Run() error: %v", err)
	}

	err = session.Close()

	if err != nil && err != io.EOF {
		return fmt.Errorf("session.Close() error: %v", err)
	}

	return nil
}
