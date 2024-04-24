package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/siruspen/logrus"

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
	server := http.Server{
		Addr:         a.address,
		Handler:      middleware.Logger(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Infof("Starting server on port %s...\n", a.address)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)

	sig := <-sigchan

	log.Info("Graceful shutdown: received ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cancel()
	server.Shutdown(ctx)
}
