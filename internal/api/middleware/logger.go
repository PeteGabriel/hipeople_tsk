package middleware

import (
	"log"
	"net/http"
)

//Log some basic info about the incoming request.
//e.g.: POST /api/image insomnia/2021.5.0
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.UserAgent())

		next.ServeHTTP(w, r)
	})
}
