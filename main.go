package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tonie-ng/go-dotenv"

	"nest/internal/api"
)

func main() {
	dotenv.Config()
	l := log.New(os.Stdout, "nest: ", log.LstdFlags)
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("We're just getting started!\n"))
	})
	server := api.New(l)
	server.Init(router)
}
