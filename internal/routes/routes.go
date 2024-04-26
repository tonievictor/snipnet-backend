package routes

import (
	"net/http"
	"os"

	"github.com/siruspen/logrus"

	"snipnet/internal/controllers"
	"snipnet/lib/services"
)

func Routes() *http.ServeMux {
	router := http.NewServeMux()
	logger := logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Up and ready to rumble!!!\n"))
	})

	users := services.User{}
	auth := controllers.NewAuthController(&users, &logger)
	router.HandleFunc("POST /signup", auth.Signup)
	return router
}
