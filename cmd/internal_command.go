package cmd

import (
	"fmt"

	homevirtcmd "github.com/nnurry/harmonia/cmd/homevirt"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/urfave/cli/v2"
)

var commandConstructorMap = map[types.InternalCommandName]types.InternalCommandConstructor{
	homevirtcmd.HOMEVIRT_COMMAND: func() types.InternalCommand { return &homevirtcmd.HomevirtCommand{} },
}

func GetCliCommand(name types.InternalCommandName) *cli.Command {
	constructor, ok := commandConstructorMap[name]
	if !ok {
		panic(fmt.Errorf("command '%v' not defined", name))
	}
	return constructor().Build()
}
