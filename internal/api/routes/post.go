package routes

import (
	"net/http"

	"learn/internal/api/handlers"
)

func RegisterPostRoutes(mux *http.ServeMux, handler *handlers.PostHandler, auth func(http.Handler) http.Handler) {
	mux.HandleFunc("GET /posts/{slug}", handler.GetPost)
	mux.Handle("POST /posts", auth(http.HandlerFunc(handler.CreatePost)))
}
