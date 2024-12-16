package main

import (
	log "log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/tonie-ng/go-dotenv"

	"snipnet/internal/api"
	"snipnet/internal/routes"
	"snipnet/lib/cache"
	"snipnet/lib/database"
	"snipnet/lib/services"
)

func main() {
	dotenv.Config("../.env")
	dbconnstr := os.Getenv("DB_CONN_STRING")

	db, err := database.Init("postgres", dbconnstr)
	if err != nil {
		log.Error("API", "Error connecting to database %v", err)
		return
	}
	defer db.Close()

	services.New(db)

	rds, err := cache.Init()
	if err != nil {
		log.Error("API", "Error connecting to redis %v", err)
		return
	}
	server := api.New(os.Getenv("PORT"))
	router := routes.Routes(rds)
	server.Init(router)
}
