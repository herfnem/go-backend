package routes

import (
	"net/http"

	"learn/internal/handlers"
	"learn/internal/middleware"
	"learn/internal/response"
)

func Setup() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /health", healthCheck)

	// Auth routes
	mux.HandleFunc("POST /auth/signup", handlers.Signup)
	mux.HandleFunc("POST /auth/login", handlers.Login)

	// Public user routes
	mux.HandleFunc("GET /users", handlers.GetUsers)
	mux.HandleFunc("GET /users/{id}", handlers.GetUser)

	// Public post routes
	mux.HandleFunc("GET /posts/{slug}", handlers.GetPost)

	// Snippet routes (public)
	mux.HandleFunc("POST /snippets", handlers.CreateSnippet)
	mux.HandleFunc("GET /s/{hash}", handlers.GetSnippet)
	mux.HandleFunc("POST /s/{hash}", handlers.GetSnippet) // POST for password-protected

	// Protected routes
	mux.Handle("GET /profile", middleware.Auth(http.HandlerFunc(handlers.GetProfile)))
	mux.Handle("POST /posts", middleware.Auth(http.HandlerFunc(handlers.CreatePost)))

	// Monitor routes (protected)
	mux.Handle("GET /monitors", middleware.Auth(http.HandlerFunc(handlers.GetMonitors)))
	mux.Handle("POST /monitors", middleware.Auth(http.HandlerFunc(handlers.CreateMonitor)))
	mux.Handle("GET /monitors/{id}", middleware.Auth(http.HandlerFunc(handlers.GetMonitor)))
	mux.Handle("DELETE /monitors/{id}", middleware.Auth(http.HandlerFunc(handlers.DeleteMonitor)))
	mux.Handle("PATCH /monitors/{id}/toggle", middleware.Auth(http.HandlerFunc(handlers.ToggleMonitor)))
	mux.Handle("GET /dashboard", middleware.Auth(http.HandlerFunc(handlers.GetDashboard)))

	// Apply global logging middleware
	return middleware.Logging(mux)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response.Success(w, "Service is healthy", map[string]string{
		"status": "healthy",
	})
}
