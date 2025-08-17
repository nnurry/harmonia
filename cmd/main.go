package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	mycli "github.com/nnurry/harmonia/cmd/cli"

	shellcmd "github.com/nnurry/harmonia/cmd/cli/shell"
	"github.com/nnurry/harmonia/internal/logger"
	"github.com/nnurry/harmonia/internal/server"
	"github.com/urfave/cli/v2"
)

func main() {
	logger.Init()

	app := &cli.App{
		Name:                 "harmonia",
		Description:          "Entrypoint of harmonia",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "start-server",
				Usage: "Start the Harmonia API server.",
				Action: func(c *cli.Context) error {
					var wg sync.WaitGroup

					osChan := make(chan os.Signal, 1)
					signal.Notify(osChan, syscall.SIGTERM, syscall.SIGINT)

					logger.Info("Starting Harmonia API server...")
					httpSrv := server.Init()

					go server.Cleanup(httpSrv, osChan, &wg)
					server.Start(httpSrv, osChan, &wg)

					wg.Wait()
					return nil
				},
			},
			mycli.GetCliCommand(shellcmd.SHELL_COMMAND),
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
