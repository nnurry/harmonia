package shell

import (
	"fmt"
	"os"

	"github.com/nnurry/harmonia/internal/homevirt/config"
	"github.com/nnurry/harmonia/internal/homevirt/processor"
	"github.com/nnurry/harmonia/internal/homevirt/service"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

const (
	SHELL_COMMAND = types.InternalCommandName("Shell command")
)

type ShellCommand struct {
	isLocal    bool
	usePrivkey bool
	cfg        config.SSH
}

func (command *ShellCommand) Description() string {
	return "Shell command entrypoint"
}

func (command *ShellCommand) Signature() string {
	return "shell"
}

func (command *ShellCommand) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "local",
			Destination: &command.isLocal,
		},
		&cli.StringFlag{
			Name:        "user",
			Value:       "root",
			Destination: &command.cfg.User,
		},
		&cli.StringFlag{
			Name:        "host",
			Value:       "localhost",
			Destination: &command.cfg.Host,
		},
		&cli.IntFlag{
			Name:        "port",
			Value:       22,
			Destination: &command.cfg.Port,
		},
		&cli.StringFlag{
			Name:        "ssh-password",
			Destination: &command.cfg.PasswordAuth.Password,
		},
		&cli.StringFlag{
			Name:        "ssh-privkey-path",
			Destination: &command.cfg.PrivateKeyAuth.PrivateKeyPath,
		},
		&cli.StringFlag{
			Name:        "ssh-passphrase",
			Destination: &command.cfg.PrivateKeyAuth.Passphrase,
		},
	}
}

func (command *ShellCommand) Subcommands() []*cli.Command {
	return []*cli.Command{}
}

func (command *ShellCommand) Handler() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if !ctx.Args().Present() {
			return fmt.Errorf("no entrypoint")
		}

		var shellProcessor processor.ShellProcessor

		if command.isLocal {
			shellProcessor = processor.NewLocalShellProcessor()
		} else {
			sshService, err := service.NewSSH(command.cfg)
			if err != nil {
				return fmt.Errorf("can't init ssh service: %v", err)
			}

			shellProcessor = processor.NewSecureShellProcessor(sshService.Client())
		}

		argsArray := ctx.Args().Slice()
		fmt.Printf("executing shell command with %v processor\n", shellProcessor.Name())
		err := shellProcessor.Execute(ctx.Context, os.Stdout, os.Stderr, argsArray[0], argsArray[1:]...)

		if err != nil {
			return fmt.Errorf("failed to execute shell command: %v", err)
		}
		return nil
	}
}

func (command *ShellCommand) Build() *cli.Command {
	command.cfg = config.SSH{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	cliCommand := utils.ConvertInternalCommandToCliCommand(command)
	return cliCommand
}
