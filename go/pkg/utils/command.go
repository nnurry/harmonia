package utils

import (
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/urfave/cli/v2"
)

func ConvertInternalCommandToCliCommand(command types.InternalCommand) *cli.Command {
	return &cli.Command{
		Name:        command.Signature(),
		Usage:       command.Description(),
		Action:      command.Handler(),
		Flags:       command.Flags(),
		Subcommands: command.Subcommands(),
	}
}
