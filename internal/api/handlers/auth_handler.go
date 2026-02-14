package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"learn/internal/api/response"
	"learn/internal/api/validator"
	"learn/internal/repository"
	"learn/internal/service"
	"learn/internal/types"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Signup godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.SignupRequest true "Signup request"
// @Success 201 {object} types.AuthResponseEnvelope
// @Failure 400 {object} types.ErrorResponseEnvelope
// @Failure 409 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req types.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, validator.FormatErrorsString(err))
		return
	}

	user, token, err := h.auth.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			response.WriteError(w, http.StatusConflict, "User already exists")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response.WriteSuccess(w, http.StatusCreated, types.AuthResponse{
		Token: token,
		User:  user.Response(),
	}, "User created successfully")
}

// Login godoc
// @Summary Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body types.LoginRequest true "Login request"
// @Success 200 {object} types.AuthResponseEnvelope
// @Failure 400 {object} types.ErrorResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validate(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, validator.FormatErrorsString(err))
		return
	}

	user, token, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	response.WriteSuccess(w, http.StatusOK, types.AuthResponse{
		Token: token,
		User:  user.Response(),
	}, "Login successful")
}
