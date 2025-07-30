package cli

import (
	"fmt"

	libvirtcmd "github.com/nnurry/harmonia/cmd/cli/libvirt"
	shellcmd "github.com/nnurry/harmonia/cmd/cli/shell"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/urfave/cli/v2"
)

var commandConstructorMap = map[types.InternalCommandName]types.InternalCommandConstructor{
	libvirtcmd.LIBVIRT_COMMAND: func() types.InternalCommand { return &libvirtcmd.LibvirtCommand{} },
	shellcmd.SHELL_COMMAND:     func() types.InternalCommand { return &shellcmd.ShellCommand{} },
}

func GetCliCommand(name types.InternalCommandName) *cli.Command {
	constructor, ok := commandConstructorMap[name]
	if !ok {
		panic(fmt.Errorf("command '%v' not defined", name))
	}
	return constructor().Build()
}
