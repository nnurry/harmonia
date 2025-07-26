package shell

import (
	"fmt"

	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

const (
	SHELL_COMMAND = types.InternalCommandName("Shell command")
)

type ShellCommand struct {
}

func (command *ShellCommand) Description() string {
	return "Shell command entrypoint"
}

func (command *ShellCommand) Signature() string {
	return "shell"
}

func (command *ShellCommand) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (command *ShellCommand) Subcommands() []*cli.Command {
	return []*cli.Command{}
}

func (command *ShellCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return fmt.Errorf("use subcommands instead")
	}
}

func (command *ShellCommand) Build() *cli.Command {
	cliCommand := utils.ConvertInternalCommandToCliCommand(command)
	return cliCommand
}
