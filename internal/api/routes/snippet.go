package routes

import (
	"net/http"

	"learn/internal/api/handlers"
)

func RegisterSnippetRoutes(mux *http.ServeMux, handler *handlers.SnippetHandler) {
	mux.HandleFunc("POST /snippets", handler.CreateSnippet)
	mux.HandleFunc("GET /s/{hash}", handler.GetSnippet)
	mux.HandleFunc("POST /s/{hash}", handler.GetSnippet)
}
