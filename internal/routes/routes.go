package routes

import "net/http"

func Routes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("We're just getting started!\n"))
	})

	return router
}
