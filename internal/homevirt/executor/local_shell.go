package executor

import (
	"os/exec"
)

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
