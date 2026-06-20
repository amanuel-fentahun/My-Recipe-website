package services

import "go-functions/internal/repository"

type AuthService struct {
	repo *repository.HasuraRepository
}

func NewAuthService(repo *repository.HasuraRepository) *AuthService {
	return &AuthService{repo: repo}
}
