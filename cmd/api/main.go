// Command api is CSMix.Matches.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tiago-bitten/CSMix.Matches/internal/shared/config"
	"github.com/tiago-bitten/CSMix.Matches/internal/shared/httpin"
	"github.com/tiago-bitten/CSMix.Matches/internal/shared/log"
)

func main() {
	logger := log.New()

	settings, err := config.Load()
	if err != nil {
		logger.Error("the service cannot start", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		httpin.NoContent(w)
	})

	handler := httpin.Chain(mux,
		httpin.RequestID(),
		httpin.Recover(logger),
		httpin.Log(logger),
	)

	server := &http.Server{
		Addr:              ":" + settings.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("listening", "port", settings.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("the server stopped", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("the shutdown did not finish cleanly", "error", err)
	}
}
