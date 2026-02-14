package middleware

import (
	"net/http"
	"strings"
)

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if isOriginAllowed(origin, allowedOrigins) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				} else if containsStar(allowedOrigins) {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PATCH,DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}
	return false
}

func containsStar(allowedOrigins []string) bool {
	for _, origin := range allowedOrigins {
		if origin == "*" {
			return true
		}
	}
	return false
}
