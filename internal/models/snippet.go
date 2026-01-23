package models

import "time"

type Snippet struct {
	ID            int64      `json:"id"`
	Hash          string     `json:"hash"`
	Content       string     `json:"content,omitempty"`
	HasPassword   bool       `json:"has_password"`
	BurnAfterRead bool       `json:"burn_after_read"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type CreateSnippetRequest struct {
	Content        string `json:"content" validate:"required,min=1,max=100000"`
	Password       string `json:"password" validate:"omitempty,min=4"`
	BurnAfterRead  bool   `json:"burn_after_read"`
	ExpiresInHours int    `json:"expires_in_hours" validate:"omitempty,min=1,max=168"`
}

type ViewSnippetRequest struct {
	Password string `json:"password"`
}

type CreateSnippetResponse struct {
	Hash      string     `json:"hash"`
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}
