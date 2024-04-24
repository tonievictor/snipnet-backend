package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"snipnet/lib/middleware"
)

type APIServer struct {
	address string
	log     *log.Logger
}

func New(l *log.Logger) *APIServer {
	apiserver := &APIServer{
		address: os.Getenv("PORT"),
		log:     l,
	}
	return apiserver
}

func (a *APIServer) Init(router *http.ServeMux) {
	server := http.Server{
		Addr:         a.address,
		Handler:      middleware.Logger(router, a.log),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	a.log.Printf("Starting server on port %s\n", a.address)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			a.log.Fatal(err)
		}
	}()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)

	sig := <-sigchan
	a.log.Println("\nGraceful shutdown: received ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cancel()
	server.Shutdown(ctx)
}
