package middleware

import (
	"net/http"
	"time"

	log "github.com/siruspen/logrus"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Infoln(r.Method, r.URL.Path, r.Body, r.UserAgent(), time.Since(start))
	})
}
