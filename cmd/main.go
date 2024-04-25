package main

import (
	"os"

	_ "github.com/lib/pq"
	log "github.com/siruspen/logrus"
	"github.com/tonie-ng/go-dotenv"

	"snipnet/internal/api"
	"snipnet/internal/routes"
	"snipnet/lib/database"
	"snipnet/lib/services"
)

func main() {
	dotenv.Config("../.env")
	db, err := database.Init("postgres", os.Getenv("DB_CONN_STRING"))
	defer db.Close()
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
		return
	}

	services.New(db)
	router := routes.Routes()
	server := api.New(os.Getenv("PORT"))
	server.Init(router)
}
