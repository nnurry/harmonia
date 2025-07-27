package shell

import (
	"fmt"
	"os"

	"github.com/nnurry/harmonia/internal/connection"
	"github.com/nnurry/harmonia/internal/processor"
	"github.com/nnurry/harmonia/pkg/types"
	"github.com/nnurry/harmonia/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

const (
	SHELL_COMMAND = types.InternalCommandName("Shell command")
)

type ShellCommand struct {
	isLocal bool
	config  connection.SSHConfig
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
			Destination: &command.config.User,
		},
		&cli.StringFlag{
			Name:        "host",
			Value:       "localhost",
			Destination: &command.config.Host,
		},
		&cli.IntFlag{
			Name:        "port",
			Value:       22,
			Destination: &command.config.Port,
		},
		&cli.StringFlag{
			Name:        "ssh-password",
			Destination: &command.config.PasswordAuth.Password,
		},
		&cli.StringFlag{
			Name:        "ssh-privkey-path",
			Destination: &command.config.PrivateKeyAuth.PrivateKeyPath,
		},
		&cli.StringFlag{
			Name:        "ssh-passphrase",
			Destination: &command.config.PrivateKeyAuth.Passphrase,
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

		var shellProcessor processor.Shell

		if command.isLocal {
			shellProcessor = processor.NewLocalShell()
		} else {
			sshConnection, err := connection.NewSSH(command.config)
			if err != nil {
				return fmt.Errorf("can't init ssh service: %v", err)
			}

			shellProcessor = processor.NewSecureShell(sshConnection)
			defer sshConnection.Cleanup()
		}

		argsArray := ctx.Args().Slice()
		log.Info().Msgf("executing shell command with %v processor\n", shellProcessor.Name())
		err := shellProcessor.Execute(ctx.Context, os.Stdout, os.Stderr, argsArray[0], argsArray[1:]...)

		if err != nil {
			return fmt.Errorf("failed to execute shell command: %v", err)
		}
		return nil
	}
}

func (command *ShellCommand) Build() *cli.Command {
	command.config = connection.SSHConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	cliCommand := utils.ConvertInternalCommandToCliCommand(command)
	return cliCommand
}
