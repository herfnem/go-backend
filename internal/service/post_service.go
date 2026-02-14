package service

import "context"

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) GetBySlug(ctx context.Context, slug string) string {
	return slug
}

func (s *PostService) Create(ctx context.Context, userID string, username string) (string, string) {
	return userID, username
}
