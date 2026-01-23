package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"learn/pkg/jwt"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
	EmailKey    contextKey = "email"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			jsonError(w, "Invalid authorization header format. Use: Bearer <token>", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			jsonError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "email", claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
