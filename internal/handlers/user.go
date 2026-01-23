package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"learn/internal/database"
	"learn/internal/models"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, email, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		jsonError(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, username, email, created_at FROM users")
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			continue
		}
		users = append(users, user)
	}

	if users == nil {
		users = []models.User{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, email, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		jsonError(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
