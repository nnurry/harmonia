package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nnurry/harmonia/internal/logger"
	"github.com/nnurry/harmonia/internal/routes"
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
	logger.Infof("serving HTTP on %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("can't start server: %v", err)
		osChan <- syscall.SIGTERM
	}
}

func Cleanup(server *http.Server, osChan chan os.Signal, wg *sync.WaitGroup) {
	sig := <-osChan
	logger.Infof("encountered OS signal %v", sig.String())

	close(osChan)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Microsecond))
	defer cancel()

	logger.Info("shutting down HTTP server")
	if err := server.Shutdown(ctx); err != nil {
		logger.Infof("called Shutdown() on HTTP server: %v", err)
	}

	logger.Info("successfully shut down HTTP server")
	wg.Done()
}
