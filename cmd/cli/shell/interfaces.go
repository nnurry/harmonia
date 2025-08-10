package shell

import (
	"context"
	"io"
)

type ShellProcessor interface {
	Name() string
	Execute(ctx context.Context, stdout io.Writer, stderr io.Writer, command string, arguments ...string) error
}
