package repository

import (
	"context"
	"errors"

	"learn/internal/models"
)

var ErrUserExists = errors.New("user already exists")

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByID(ctx context.Context, id string) (models.User, error)
	List(ctx context.Context) ([]models.User, error)
}
