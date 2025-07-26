package processor

import (
	"context"
	"io"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

type ShellProcessor interface {
	Execute(context.Context, string, ...string)
}

type LocalShellProcessor struct {
}

func (processor *LocalShellProcessor) Execute(
	ctx context.Context,
	stdout, stderr io.Writer,
	command string, args ...string,
) error {
	cmd := exec.CommandContext(ctx, command, args...)
	err := cmd.Start()

	if err != nil {
		return err
	}
	return nil
}

type SecureShellProcessor struct {
	client *ssh.Client
}

func NewSecureShellProcessor(client *ssh.Client) *SecureShellProcessor {
	return &SecureShellProcessor{client: client}
}

func (processor *SecureShellProcessor) Execute(
	ctx context.Context,
	stdout, stderr io.Writer,
	command string, args ...string,
) error {
	session, err := processor.client.NewSession()
	if err != nil {
		return err
	}

	defer session.Close()

	session.Stdout = stdout
	session.Stderr = stderr

	commandParts := []string{command}
	commandParts = append(commandParts, args...)

	command = strings.Join(commandParts, " ")

	err = session.Start(command)
	if err != nil {
		return err
	}

	<-ctx.Done()
	err = session.Signal(ssh.SIGTERM)

	if err != nil {
		return err
	}

	return nil
}
