package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/siruspen/logrus"
	log "github.com/siruspen/logrus"

	"snipnet/internal/utils"
	"snipnet/lib/types"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Infoln(r.Method, r.URL.Path, r.UserAgent(), time.Since(start))
	})
}

func IsAuthenticated(next func(http.ResponseWriter, *http.Request), log *logrus.Logger, cache *redis.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			utils.WriteErr(w, http.StatusUnauthorized, "You are not logged in", errors.New("Session id not found"), log)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			utils.WriteErr(w, http.StatusUnauthorized, "You are not logged in", errors.New("Session id not found"), log)
			return
		}

		val, err := cache.Get(context.Background(), token).Result()
		if err == redis.Nil || err != nil {
			utils.WriteErr(w, http.StatusUnauthorized, "Invalid session token", err, log)
			return
		}
		var session types.Session

		err = json.Unmarshal([]byte(val), &session)
		if err != nil {
			utils.WriteErr(w, http.StatusInternalServerError, "An error occured while validating token ", err, log)
			return
		}

		ctx := context.WithValue(r.Context(), types.AuthSession, session)
		req := r.WithContext(ctx)

		next(w, req)
	}
}
