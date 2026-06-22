package actions

import (
	"go-functions/internal/repository"
	"go-functions/services"
	"os"
)

var (
	HasuraRepo  = repository.NewHasuraRepository()
	authService = services.NewAuthService(
		HasuraRepo,
		os.Getenv("JWT_SECRET"),
		os.Getenv("VERIFICATION_CODE_TTL"), // e.g., "15m" or "1h"
	)
	verifyService = services.NewVerificationService(HasuraRepo)
)
