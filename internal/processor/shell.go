package processor

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/nnurry/harmonia/internal/connection"
)

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
	client *connection.SSH
}

func NewSecureShell(client *connection.SSH) *SecureShell {
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
	session, err := processor.client.Session()
	if err != nil {
		return fmt.Errorf("client.Session() error: %v", err)
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
