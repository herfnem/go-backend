package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"learn/internal/database"
	"learn/internal/models"
	"learn/pkg/jwt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		jsonError(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		jsonError(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		req.Username, req.Email, string(hashedPassword),
	)
	if err != nil {
		jsonError(w, "User already exists or database error", http.StatusConflict)
		return
	}

	userID, _ := result.LastInsertId()

	token, err := jwt.GenerateToken(userID, req.Username, req.Email)
	if err != nil {
		jsonError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token:   token,
		Message: "User created successfully",
		User:    user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		jsonError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, email, password, created_at FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err == sql.ErrNoRows {
		jsonError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		jsonError(w, "Database error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		jsonError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		jsonError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	user.Password = "" // Clear password before sending response

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token:   token,
		Message: "Login successful",
		User:    user,
	})
}

func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
