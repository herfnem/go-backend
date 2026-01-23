package handlers

import (
	"database/sql"
	"net/http"

	"learn/internal/database"
	"learn/internal/middleware"
	"learn/internal/models"
	"learn/internal/response"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, email, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err == sql.ErrNoRows {
		response.NotFound(w, "User not found")
		return
	}
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}

	response.Success(w, "User retrieved successfully", user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, username, email, created_at FROM users")
	if err != nil {
		response.InternalError(w, "Database error")
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

	response.Success(w, "Users retrieved successfully", users)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	response.Success(w, "Profile retrieved successfully", user)
}
