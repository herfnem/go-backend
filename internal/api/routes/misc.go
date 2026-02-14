package routes

import (
	"net/http"

	"learn/internal/api/handlers"
)

func RegisterMiscRoutes(mux *http.ServeMux, handler *handlers.MiscHandler) {
	mux.HandleFunc("GET /{$}", handler.Home)
	mux.HandleFunc("GET /health", handler.Health)
}
