package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"learn/internal/database"
	"learn/internal/models"
	"learn/internal/response"
	"learn/pkg/jwt"
)

type contextKey string

const (
	UserKey contextKey = "user"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(w, "Authorization header required")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(w, "Invalid authorization header format. Use: Bearer <token>")
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(w, "Invalid or expired token")
			return
		}

		// Fetch fresh user data from database
		var user models.User
		err = database.DB.QueryRow(
			"SELECT id, username, email, created_at FROM users WHERE id = ?",
			claims.UserID,
		).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

		if err == sql.ErrNoRows {
			response.Unauthorized(w, "User not found")
			return
		}
		if err != nil {
			response.InternalError(w, "Database error")
			return
		}

		// Inject full user object into context
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext helper to retrieve user from context
func GetUserFromContext(r *http.Request) (models.User, bool) {
	user, ok := r.Context().Value(UserKey).(models.User)
	return user, ok
}
