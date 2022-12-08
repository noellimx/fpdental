package middlewares

import (
	"log"
	"net/http"
)

func RestJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Middleware Restful called.")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
