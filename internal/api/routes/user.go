package routes

import (
	"net/http"

	"learn/internal/api/handlers"
)

func RegisterUserRoutes(mux *http.ServeMux, handler *handlers.UserHandler, auth func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /users", handler.GetUsers)
	mux.HandleFunc("GET /users/{id}", handler.GetUser)
	mux.Handle("GET /profile", auth(http.HandlerFunc(handler.GetProfile)))
}
