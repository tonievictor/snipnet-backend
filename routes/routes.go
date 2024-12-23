package routes

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"

	"snipnet/controllers"
	"snipnet/controllers/middleware"
	"snipnet/services"
)

func Routes(rds *redis.Client) http.Handler {
	mux := http.NewServeMux()
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

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

	// swagger
	handleFunc("GET /swagger/*", httpSwagger.Handler())
	handleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("Up and ready to rumble!!!\n"))
	})

	users := services.User{}
	auth_controller := controllers.NewAuthController(&users, &oauthConfig, logger, rds)
	handleFunc("POST /signin", auth_controller.GitHubOauth)
	handleFunc("POST /signout", middleware.IsAuthenticated(auth_controller.Signout, logger, rds))

	snippets := services.Snippet{}
	snippet_controller := controllers.NewSnippetController(&snippets, logger, rds)
	handleFunc("GET /snippets/{id}", snippet_controller.GetSnippetByID)
	handleFunc("GET /snippets", snippet_controller.GetAllSnippets)
	handleFunc("POST /snippets", middleware.IsAuthenticated(snippet_controller.CreateSnippet, logger, rds))
	handleFunc("DELETE /snippets/{id}", middleware.IsAuthenticated(snippet_controller.DeleteSnippet, logger, rds))
	handleFunc("PUT /snippets/{id}", middleware.IsAuthenticated(snippet_controller.UpdateSnippetMulti, logger, rds))
	handleFunc("PATCH /snippets/{id}", middleware.IsAuthenticated(snippet_controller.UpdateSnippetOne, logger, rds))

	user_controller := controllers.NewUserController(&users, logger, rds)
	handleFunc("GET /users/{id}", middleware.IsAuthenticated(user_controller.GetUserByID, logger, rds))
	handleFunc("GET /users/{id}/snippets", middleware.IsAuthenticated(snippet_controller.GetAllUserSnippets, logger, rds))

	// add cors
	handler := otelhttp.NewHandler(mux, "/")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
		Debug:            true,
	})
	router := c.Handler(handler)

	return router
}
