package routes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"

	"snipnet/internal/controllers"
	"snipnet/lib/middleware"
	"snipnet/lib/services"
)

func Routes(rds *redis.Client) *http.ServeMux {
	router := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	oauthConfig := oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: os.Getenv("GH_REDIRECT_URL"),
		Scopes:      []string{"user"},
	}

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Up and ready to rumble!!!\n"))
	})

	users := services.User{}
	auth_controller := controllers.NewAuthController(&users, &oauthConfig, logger, rds)
	router.HandleFunc("GET /signin", auth_controller.GitHubOauth)
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
	return router
}
