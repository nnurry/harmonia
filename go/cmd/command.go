package cmd

import (
	"fmt"

	libvirt_cmd "github.com/nnurry/harmonia/cmd/libvirt"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
)

var commandConstructorMap = map[types.InternalCommandName]types.InternalCommandConstructor{
	libvirt_cmd.LIBVIRT_COMMAND: func() types.InternalCommand { return &libvirt_cmd.LibvirtCommand{} },
}

func GetCliCommand(name types.InternalCommandName) *cli.Command {
	constructor, ok := commandConstructorMap[name]
	if !ok {
		panic(fmt.Errorf("command '%v' not defined", name))
	}
	command := constructor()
	return utils.ConvertInternalCommandToCliCommand(command)
}
