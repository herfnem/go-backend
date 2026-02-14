package routes

import (
	"net/http"

	"learn/internal/api/handlers"
)

func RegisterAuthRoutes(mux *http.ServeMux, handler *handlers.AuthHandler) {
	mux.HandleFunc("POST /auth/signup", handler.Signup)
	mux.HandleFunc("POST /auth/login", handler.Login)
}
