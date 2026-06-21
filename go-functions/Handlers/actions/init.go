package actions

import (
	"go-functions/internal/repository"
	"go-functions/services"
	"os"
)

var (
	HasuraRepo    = repository.NewHasuraRepository()
	authService   = services.NewAuthService(HasuraRepo, os.Getenv("JWT_SECRET"))
	verifyService = services.NewVerificationService(HasuraRepo)
)
