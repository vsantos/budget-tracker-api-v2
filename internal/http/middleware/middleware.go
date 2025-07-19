package middleware

import (
	"log"
	"net/http"
)

// InjectHeaders acts as a middleware between routers to inject common response headers
func InjectHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
