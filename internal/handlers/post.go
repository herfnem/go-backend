package handlers

import (
	"net/http"

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
	userID := r.Context().Value("user_id").(int64)

	response.Created(w, "Post created successfully", map[string]any{
		"user_id": userID,
	})
}
