package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"learn/internal/database"
	"learn/internal/models"
	"learn/internal/response"
	"learn/pkg/jwt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.BadRequest(w, "Username, email, and password are required")
		return
	}

	if len(req.Password) < 6 {
		response.BadRequest(w, "Password must be at least 6 characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(w, "Failed to process password")
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		req.Username, req.Email, string(hashedPassword),
	)
	if err != nil {
		response.Conflict(w, "User already exists")
		return
	}

	userID, _ := result.LastInsertId()

	token, err := jwt.GenerateToken(userID, req.Username, req.Email)
	if err != nil {
		response.InternalError(w, "Failed to generate token")
		return
	}

	user := models.User{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	response.Created(w, "User created successfully", map[string]any{
		"token": token,
		"user":  user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.BadRequest(w, "Email and password are required")
		return
	}

	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, email, password, created_at FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err == sql.ErrNoRows {
		response.Unauthorized(w, "Invalid email or password")
		return
	}
	if err != nil {
		response.InternalError(w, "Database error")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		response.Unauthorized(w, "Invalid email or password")
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		response.InternalError(w, "Failed to generate token")
		return
	}

	user.Password = ""

	response.Success(w, "Login successful", map[string]any{
		"token": token,
		"user":  user,
	})
}
