package routes

import (
	"net/http"

	"learn/internal/handlers"
	"learn/internal/middleware"
	"learn/internal/response"
)

func Setup() http.Handler {
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /health", healthCheck)

	// Auth routes
	mux.HandleFunc("POST /auth/signup", handlers.Signup)
	mux.HandleFunc("POST /auth/login", handlers.Login)

	// Public post routes
	mux.HandleFunc("GET /posts/{slug}", handlers.GetPost)

	// Public user routes
	mux.HandleFunc("GET /users", handlers.GetUsers)
	mux.HandleFunc("GET /users/{id}", handlers.GetUser)

	// Protected routes (require authentication)
	mux.Handle("GET /profile", middleware.Auth(http.HandlerFunc(handlers.GetProfile)))
	mux.Handle("POST /posts", middleware.Auth(http.HandlerFunc(handlers.CreatePost)))

	// Apply global logging middleware
	return middleware.Logging(mux)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response.Success(w, "Service is healthy", map[string]string{
		"status": "healthy",
	})
}
