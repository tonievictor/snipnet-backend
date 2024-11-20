package main

import (
	log "log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/tonie-ng/go-dotenv"

	"snipnet/internal/api"
	"snipnet/internal/routes"
	"snipnet/lib/database"
	"snipnet/lib/services"
)

func main() {
	dotenv.Config("../.env")
	dbconnstr := os.Getenv("DB_CONN_STRING")

	db, err := database.Init("postgres", dbconnstr)

	defer db.Close()
	if err != nil {
		log.Error("API", "Error connecting to database %v", err)
		return
	}

	services.New(db)
	router := routes.Routes()
	server := api.New(os.Getenv("PORT"))
	server.Init(router)
}
