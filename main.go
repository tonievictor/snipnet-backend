package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/tonie-ng/go-dotenv"

	"snipnet/internal/api"
	"snipnet/lib/database"
)

func main() {
	dotenv.Config()
	l := log.New(os.Stdout, "nest: ", log.LstdFlags)
	db := database.New("postgres", os.Getenv("DB_CONN_STRING"), l)
	db.Init()
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("We're just getting started!\n"))
	})

	server := api.New(l)
	server.Init(router)
}
