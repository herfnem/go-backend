package models

import "time"

type Snippet struct {
	ID            string     `json:"id"`
	Hash          string     `json:"hash"`
	Content       string     `json:"content"`
	PasswordHash  *string    `json:"-"`
	BurnAfterRead bool       `json:"burn_after_read"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type SnippetResponse struct {
	ID            string     `json:"id"`
	Hash          string     `json:"hash"`
	Content       string     `json:"content,omitempty"`
	HasPassword   bool       `json:"has_password"`
	BurnAfterRead bool       `json:"burn_after_read"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (s Snippet) Response() SnippetResponse {
	hasPassword := s.PasswordHash != nil && *s.PasswordHash != ""
	return SnippetResponse{
		ID:            s.ID,
		Hash:          s.Hash,
		Content:       s.Content,
		HasPassword:   hasPassword,
		BurnAfterRead: s.BurnAfterRead,
		ExpiresAt:     s.ExpiresAt,
		CreatedAt:     s.CreatedAt,
	}
}
