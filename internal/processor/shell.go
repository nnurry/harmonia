package processor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/nnurry/harmonia/internal/connection"
	"github.com/rs/zerolog/log"
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
	cmdBytesBuffer := bytes.NewBufferString(command)
	for _, arg := range args {
		cmdBytesBuffer.WriteString(" ")
		cmdBytesBuffer.WriteString(arg)
	}

	log.Info().Msgf("executing command '%v' locally", cmdBytesBuffer.String())
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

	cmdBytesBuffer := bytes.NewBufferString(command)
	for _, arg := range args {
		cmdBytesBuffer.WriteString(" ")
		cmdBytesBuffer.WriteString(arg)
	}

	command = cmdBytesBuffer.String()

	session.Stdout = stdout
	session.Stderr = stderr

	log.Info().Msgf("executing command '%v' via SSH", command)

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
