package service

import (
	"context"

	"learn/internal/models"
	"learn/internal/repository"
)

type UserService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetByID(ctx context.Context, id string) (models.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.users.List(ctx)
}
