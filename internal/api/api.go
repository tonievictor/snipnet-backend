package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"

	"snipnet/lib/middleware"
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

func (a *APIServer) Init(router *http.ServeMux) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
		Debug:            true,
	})
	handler := c.Handler(router)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", a.address),
		Handler:      middleware.Logger(handler),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	slog.Info("API", slog.String("Server runnning on", a.address))
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("API", slog.String("error", err.Error()))
			return
		}
	}()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)
	sig := <-sigchan

	slog.Info("API", slog.String("Graceful shutdown: received %v\n", sig.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cancel()
	server.Shutdown(ctx)
}
