package routes

import (
	"net/http"
	"os"

	"github.com/siruspen/logrus"

	"snipnet/internal/controllers"
	"snipnet/lib/cache"
	"snipnet/lib/middleware"
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
	rds := cache.Init(os.Getenv("REDIS_CLIENT"))

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Up and ready to rumble!!!\n"))
	})

	users := services.User{}
	auth_controller := controllers.NewAuthController(&users, &logger, rds)
	router.HandleFunc("POST /signup", auth_controller.Signup)
	router.HandleFunc("POST /signin", auth_controller.Signin)
	router.HandleFunc("POST /signout", middleware.IsAuthenticated(auth_controller.Signout, &logger, rds))

	snippets := services.Snippet{}
	snippet_controller := controllers.NewSnippetController(&snippets, &logger, rds)
	router.HandleFunc("GET /snippets/{id}", middleware.IsAuthenticated(snippet_controller.GetSnippetByID, &logger, rds))
	router.HandleFunc("POST /snippets", middleware.IsAuthenticated(snippet_controller.CreateSnippet, &logger, rds))
	return router
}
