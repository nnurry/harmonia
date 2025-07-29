package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nnurry/harmonia/internal/routes"
	"github.com/rs/zerolog/log"
)

func Init() {
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGTERM, syscall.SIGINT)

	mux := routes.SetupMux()
	httpSrv := http.Server{
		Addr:    ":15000",
		Handler: mux,
	}

	go Cleanup(&httpSrv, osChan)
}

func Cleanup(server *http.Server, osChan chan os.Signal) {
	log.Info().Msgf("serving HTTP on :%v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Msgf("can't start server: %v", err)
		osChan <- syscall.SIGTERM
	}
}
