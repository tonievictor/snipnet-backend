package main

import (
	log "log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/tonie-ng/go-dotenv"

	"snipnet/docs"
	"snipnet/api"
	"snipnet/routes"

	"snipnet/database/cache"
	"snipnet/database"
	"snipnet/services"
)

func main() {
	docs.SwaggerInfo.Title = "Snipnet API"
	docs.SwaggerInfo.Description = "API for Snipnet, a code snippet storage and sharing platform."
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = os.Getenv("API_HOST") // e.g localhost:8080
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	dotenv.Config(".env")
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
