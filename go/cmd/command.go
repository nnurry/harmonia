package cmd

import "github.com/urfave/cli/v2"

type Command interface {
	Description() string
	Signature() string
	Handler() func(ctx *cli.Context) error
}
