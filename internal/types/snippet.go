package types

import (
	"time"

	"learn/internal/models"
)

type SnippetCreateRequest struct {
	Content        string `json:"content" validate:"required,min=1,max=100000" example:"This is my secret message"`
	Password       string `json:"password" validate:"omitempty,min=4" example:"optional-password"`
	BurnAfterRead  bool   `json:"burn_after_read" example:"true"`
	ExpiresInHours int    `json:"expires_in_hours" validate:"omitempty,min=1,max=168" example:"24"`
}

type SnippetViewRequest struct {
	Password string `json:"password" example:"optional-password"`
}

type SnippetCreateResponse struct {
	Hash      string     `json:"hash"`
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type SnippetCreateResponseEnvelope struct {
	Success bool                  `json:"success"`
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    SnippetCreateResponse `json:"data"`
}

type SnippetResponseEnvelope struct {
	Success bool                   `json:"success"`
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    models.SnippetResponse `json:"data"`
}

type SnippetPasswordRequiredResponse struct {
	PasswordRequired bool   `json:"password_required"`
	Hash             string `json:"hash"`
}

type SnippetPasswordRequiredEnvelope struct {
	Success bool                            `json:"success"`
	Status  int                             `json:"status"`
	Message string                          `json:"message"`
	Data    SnippetPasswordRequiredResponse `json:"data"`
}
