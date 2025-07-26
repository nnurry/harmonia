package processor

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

type ShellProcessor interface {
	Name() string
	Execute(context.Context, io.Writer, io.Writer, string, ...string) error
}

type LocalShellProcessor struct {
}

func NewLocalShellProcessor() *LocalShellProcessor {
	return &LocalShellProcessor{}
}

func (processor *LocalShellProcessor) Name() string {
	return "local-shell"
}

func (processor *LocalShellProcessor) Execute(
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

type SecureShellProcessor struct {
	client *ssh.Client
}

func NewSecureShellProcessor(client *ssh.Client) *SecureShellProcessor {
	return &SecureShellProcessor{client: client}
}

func (processor *SecureShellProcessor) Name() string {
	return "secure-shell"
}

func (processor *SecureShellProcessor) Execute(
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
