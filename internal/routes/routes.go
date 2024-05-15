package routes

import (
	"log/slog"
	"net/http"
	"os"

	"snipnet/internal/controllers"
	"snipnet/lib/cache"
	"snipnet/lib/middleware"
	"snipnet/lib/services"
)

func Routes() *http.ServeMux {
	router := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	rds := cache.Init(os.Getenv("REDIS_CLIENT"))

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Up and ready to rumble!!!\n"))
	})

	users := services.User{}
	auth_controller := controllers.NewAuthController(&users, logger, rds)
	router.HandleFunc("POST /signup", auth_controller.Signup)
	router.HandleFunc("POST /signin", auth_controller.Signin)
	router.HandleFunc("POST /signout", middleware.IsAuthenticated(auth_controller.Signout, logger, rds))

	snippets := services.Snippet{}
	snippet_controller := controllers.NewSnippetController(&snippets, logger, rds)
	router.HandleFunc("GET /snippets/{id}", snippet_controller.GetSnippetByID)
	router.HandleFunc("GET /snippets", snippet_controller.GetAllSnippets)
	router.HandleFunc("POST /snippets", middleware.IsAuthenticated(snippet_controller.CreateSnippet, logger, rds))
	router.HandleFunc("DELETE /snippets/{id}", middleware.IsAuthenticated(snippet_controller.DeleteSnippet, logger, rds))
	router.HandleFunc("PUT /snippets/{id}", middleware.IsAuthenticated(snippet_controller.UpdateSnippetMulti, logger, rds))
	router.HandleFunc("PATCH /snippets/{id}", middleware.IsAuthenticated(snippet_controller.UpdateSnippetOne, logger, rds))

	user_controller := controllers.NewUserController(&users, logger, rds)
	router.HandleFunc("GET /users", middleware.IsAuthenticated(user_controller.GetUsers, logger, rds))
	router.HandleFunc("GET /users/{id}", middleware.IsAuthenticated(user_controller.GetUserByID, logger, rds))
	router.HandleFunc("GET /users/{id}/snippets", middleware.IsAuthenticated(snippet_controller.GetAllUserSnippets, logger, rds))
	router.HandleFunc("GET /users/snippets", middleware.IsAuthenticated(snippet_controller.GetAllCurrentUserSnippets, logger, rds))
	router.HandleFunc("PATCH /users/{id}", middleware.IsAuthenticated(user_controller.UpdateUserOne, logger, rds))
	router.HandleFunc("PUT /users/{id}", middleware.IsAuthenticated(user_controller.UpdateUserMulti, logger, rds))
	router.HandleFunc("DELETE /users/{id}", middleware.IsAuthenticated(user_controller.DeleteUser, logger, rds))
	return router
}
