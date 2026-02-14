package types

import "learn/internal/models"

type UserResponseEnvelope struct {
	Success bool                `json:"success"`
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    models.UserResponse `json:"data"`
}

type UsersResponseEnvelope struct {
	Success bool                  `json:"success"`
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    []models.UserResponse `json:"data"`
}
