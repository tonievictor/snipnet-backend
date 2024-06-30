package main

import (
	"fmt"
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
	dbconnstr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
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
