package handlers

import (
	"net/http"

	"learn/internal/middleware"
	"learn/internal/response"
)

func Home(w http.ResponseWriter, r *http.Request) {
	response.Success(w, "Welcome to the API", map[string]string{
		"version": "1.0.0",
	})
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	response.Success(w, "Post retrieved successfully", map[string]string{
		"slug": slug,
	})
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	response.Created(w, "Post created successfully", map[string]any{
		"user_id":  user.ID,
		"username": user.Username,
	})
}
