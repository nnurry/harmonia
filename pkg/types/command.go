package types

import "github.com/urfave/cli/v2"

type InternalCommandName string
type InternalCommandCtxKey string
type InternalCommandConstructor func() InternalCommand

type InternalCommand interface {
	Description() string
	Signature() string
	Handler() func(ctx *cli.Context) error
	Flags() []cli.Flag
	Subcommands() []*cli.Command
	Build() *cli.Command
}
