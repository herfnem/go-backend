package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"learn/internal/api/response"
	"learn/internal/models"
	"learn/internal/repository"
	"learn/pkg/jwt"
)

type contextKey string

const (
	UserKey contextKey = "user"
)

func Auth(users repository.UserRepository, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.WriteError(w, http.StatusUnauthorized, "Invalid authorization header format. Use: Bearer <token>")
				return
			}

			tokenString := parts[1]
			claims, err := jwt.ValidateToken(jwtSecret, tokenString)
			if err != nil {
				response.WriteError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			user, err := users.GetByID(r.Context(), claims.UserID)
			if err == sql.ErrNoRows {
				response.WriteError(w, http.StatusUnauthorized, "User not found")
				return
			}
			if err != nil {
				response.WriteError(w, http.StatusInternalServerError, "Database error")
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(r *http.Request) (models.User, bool) {
	user, ok := r.Context().Value(UserKey).(models.User)
	return user, ok
}
