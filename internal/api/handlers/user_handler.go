package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"learn/internal/api/middleware"
	"learn/internal/api/response"
	"learn/internal/models"
	"learn/internal/service"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

// GetUsers godoc
// @Summary List all users
// @Tags users
// @Produce json
// @Success 200 {object} types.UsersResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	result := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, user.Response())
	}

	response.WriteSuccess(w, http.StatusOK, result, "Users retrieved successfully")
}

// GetUser godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} types.UserResponseEnvelope
// @Failure 400 {object} types.ErrorResponseEnvelope
// @Failure 404 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			response.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	response.WriteSuccess(w, http.StatusOK, user.Response(), "User retrieved successfully")
}

// GetProfile godoc
// @Summary Get current user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} types.UserResponseEnvelope
// @Failure 401 {object} types.ErrorResponseEnvelope
// @Failure 500 {object} types.ErrorResponseEnvelope
// @Router /profile [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	response.WriteSuccess(w, http.StatusOK, user.Response(), "Profile retrieved successfully")
}
