package services

import (
	"context"
	"errors"
	"time"

	"go-functions/internal/repository"
)

var HasuraRepo = repository.NewHasuraRepository()

func VerifyEmailService(email string, old_code string, incoming_code string, expired_at time.Time, ctx context.Context) (bool, error) {

	if time.Now().After(expired_at) {
		return false, errors.New("verification code has expired")
	}

	// check the code
	if old_code != incoming_code {
		return false, errors.New("invalid verification code")
	}

	// update the user table isVerified to true
	if err := HasuraRepo.MarkEmailVerified(ctx, email); err != nil {
		return false, err
	}

	return true, nil
}
