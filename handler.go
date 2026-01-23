package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserRequest struct {
	Username string `json:"Username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

type UserResponse struct {
	ID      int64  `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"Username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

func handleHome(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Welcome to home route."))
}

func handlePostView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	fmt.Fprintf(w, "Displaying Post: %s", slug)
}

func handlePostCreate(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Post created successfully."))
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", 400)
		return
	}

	// Insert into Database
	result, err := DB.Exec("INSERT INTO users (username, email) VALUES (?, ?)", req.Username, req.Email)
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}

	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":      id,
		"message": "User saved to database!",
	})
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var u User
	err := DB.QueryRow("SELECT id, username, email FROM users WHERE id = ?", id).Scan(&u.ID, &u.Username, &u.Email)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT username, email FROM users")
	if err != nil {
		http.Error(w, "Database error", 500)
		return
	}
	defer rows.Close() // ALWAYS close your rows to avoid memory leaks

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Username, &u.Email); err != nil {
			continue
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

}
