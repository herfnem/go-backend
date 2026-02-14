package repository

import (
	"context"
	"time"

	"learn/internal/models"
)

type SnippetRepository interface {
	Create(ctx context.Context, snippet models.Snippet) (models.Snippet, error)
	GetByHash(ctx context.Context, hash string) (models.Snippet, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context, now time.Time) (int64, error)
}
