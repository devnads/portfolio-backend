package middleware

import (
	"net/http"
	"os"
)

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminToken := os.Getenv("ADMIN_TOKEN")
		clientToken := r.Header.Get("X-Admin-Token")

		if clientToken == ""  || clientToken != adminToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w,r)
	})
}