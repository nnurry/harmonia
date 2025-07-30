package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nnurry/harmonia/internal/routes"
	"github.com/rs/zerolog/log"
)

func Init() *http.Server {
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGTERM, syscall.SIGINT)

	mux := routes.SetupMux()
	httpSrv := http.Server{
		Addr:    ":15000",
		Handler: mux,
	}

	return &httpSrv
}

func Start(server *http.Server, osChan chan os.Signal, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info().Msgf("serving HTTP on %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Msgf("can't start server: %v", err)
		osChan <- syscall.SIGTERM
	}
}

func Cleanup(server *http.Server, osChan chan os.Signal, wg *sync.WaitGroup) {
	sig := <-osChan
	log.Info().Msgf("encountered OS signal %v", sig.String())

	close(osChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Microsecond))
	defer cancel()

	log.Info().Msg("shutting down HTTP server")
	if err := server.Shutdown(ctx); err != nil {
		log.Info().Msgf("called Shutdown() on HTTP server: %v", err)
	}

	log.Info().Msg("successfully shut down HTTP server")
	wg.Done()
}
