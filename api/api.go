package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"snipnet/controllers/middleware"
)

type APIServer struct {
	address string
}

func New(port string) *APIServer {
	apiserver := &APIServer{
		address: port,
	}
	return apiserver
}

func (a *APIServer) Init(router http.Handler) {
	server := http.Server{
		Addr:         fmt.Sprintf(":%s", a.address),
		Handler:      middleware.Logger(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	slog.Info("API", slog.String("Server runnning on", a.address))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	otelShutDown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}

	defer func() {
		err = errors.Join(err, otelShutDown(context.Background()))
	}()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err = <-serverErr:
		slog.Error("API", slog.String("error", err.Error()))
		return
	case <-ctx.Done():
		slog.Info("API", slog.String("SHUTDOWN", "Graceful shutdown: received\n"))
		stop()
	}
	err = server.Shutdown(context.Background())
	return
}
