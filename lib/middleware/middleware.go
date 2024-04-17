package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler, l *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		l.Println(r.Method, r.URL.Path, time.Since(start))
	})
}
