package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"learn/internal/models"
	"learn/internal/repository"
	"learn/pkg/jwt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	users     repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthService(users repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{users: users, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (models.User, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, "", err
	}

	user, err := s.users.Create(ctx, models.User{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return models.User{}, "", err
	}

	token, err := jwt.GenerateToken(s.jwtSecret, s.jwtExpiry, user.ID, user.Username, user.Email)
	if err != nil {
		return models.User{}, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (models.User, string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, "", ErrInvalidCredentials
		}
		return models.User{}, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.User{}, "", ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(s.jwtSecret, s.jwtExpiry, user.ID, user.Username, user.Email)
	if err != nil {
		return models.User{}, "", err
	}

	return user, token, nil
}
