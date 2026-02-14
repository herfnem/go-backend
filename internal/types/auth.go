package types

import "learn/internal/models"

type SignupRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum" example:"john"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"secret123"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required" example:"secret123"`
}

type AuthResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

type AuthResponseEnvelope struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    AuthResponse `json:"data"`
}
