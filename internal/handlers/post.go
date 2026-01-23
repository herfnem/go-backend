package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to the API",
		"version": "1.0.0",
	})
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"slug":    slug,
		"message": fmt.Sprintf("Displaying post: %s", slug),
	})
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// This handler requires authentication (added via middleware in routes)
	userID := r.Context().Value("user_id").(int64)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Post created successfully",
		"user_id": userID,
	})
}
