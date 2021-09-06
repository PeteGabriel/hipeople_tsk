package middleware

import (
	"log"
	"net/http"
)

//Log some basic info about the incoming request.
//e.g.: POST /api/image <user-agent>
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.UserAgent())

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
