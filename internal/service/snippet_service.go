package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"learn/internal/models"
	"learn/internal/repository"
)

var (
	ErrSnippetNotFound         = errors.New("snippet not found")
	ErrSnippetPasswordRequired = errors.New("password required")
	ErrSnippetInvalidPassword  = errors.New("invalid password")
)

type SnippetService struct {
	snippets repository.SnippetRepository
}

func NewSnippetService(snippets repository.SnippetRepository) *SnippetService {
	return &SnippetService{snippets: snippets}
}

func (s *SnippetService) Create(ctx context.Context, content, password string, burnAfterRead bool, expiresInHours int) (models.Snippet, error) {
	hash, err := generateHash(8)
	if err != nil {
		return models.Snippet{}, err
	}

	var hashedPassword *string
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return models.Snippet{}, err
		}
		hashedStr := string(hashed)
		hashedPassword = &hashedStr
	}

	var expiresAt *time.Time
	if expiresInHours > 0 {
		exp := time.Now().Add(time.Duration(expiresInHours) * time.Hour)
		expiresAt = &exp
	} else {
		exp := time.Now().Add(24 * time.Hour)
		expiresAt = &exp
	}

	return s.snippets.Create(ctx, models.Snippet{
		ID:            uuid.NewString(),
		Hash:          hash,
		Content:       content,
		PasswordHash:  hashedPassword,
		BurnAfterRead: burnAfterRead,
		ExpiresAt:     expiresAt,
	})
}

func (s *SnippetService) Get(ctx context.Context, hash, password string) (models.Snippet, error) {
	snippet, err := s.snippets.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Snippet{}, ErrSnippetNotFound
		}
		return models.Snippet{}, err
	}

	if snippet.ExpiresAt != nil && time.Now().After(*snippet.ExpiresAt) {
		_ = s.snippets.DeleteByID(ctx, snippet.ID)
		return models.Snippet{}, ErrSnippetNotFound
	}

	if snippet.PasswordHash != nil && *snippet.PasswordHash != "" {
		if password == "" {
			return models.Snippet{}, ErrSnippetPasswordRequired
		}
		if err := bcrypt.CompareHashAndPassword([]byte(*snippet.PasswordHash), []byte(password)); err != nil {
			return models.Snippet{}, ErrSnippetInvalidPassword
		}
	}

	if snippet.BurnAfterRead {
		_ = s.snippets.DeleteByID(ctx, snippet.ID)
	}

	return snippet, nil
}

func (s *SnippetService) DeleteExpired(ctx context.Context) (int64, error) {
	return s.snippets.DeleteExpired(ctx, time.Now())
}

func generateHash(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}
