package cmd

import "github.com/urfave/cli/v2"

type DefineLibvirtDomainCommand struct {
}

func (command *DefineLibvirtDomainCommand) Description() string {
	return "Define a domain in Libvirt"
}

func (command *DefineLibvirtDomainCommand) Signature() string {
	return "define-domain"
}

func (command *DefineLibvirtDomainCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		return nil
	}
}

func (command DefineLibvirtDomainCommand) BuildCliCommand() *cli.Command {
	return &cli.Command{
		Name:   command.Signature(),
		Usage:  command.Description(),
		Action: command.Handler(),
	}
}
