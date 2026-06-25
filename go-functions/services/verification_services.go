package services

import (
	"context"
	"errors"
	"go-functions/internal/repository"
	"go-functions/internal/response"
	"net/http"
	"time"
)

type VerificationService struct {
	repo *repository.HasuraRepository
}

func NewVerificationService(repo *repository.HasuraRepository) *VerificationService {
	return &VerificationService{repo: repo}
}

func (s *VerificationService) VerifyEmail(ctx context.Context, email, old_code, incoming_code string, expiredAt time.Time) error {

	if time.Now().After(expiredAt) {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "Verification code has expired. Please request a new one.",
			RawError:   errors.New("time window elapsed"),
		}
	}

	// check the code
	if old_code != incoming_code {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "Invalid verification code provided.",
			RawError:   errors.New("code verification mismatch"),
		}
	}

	return nil
}

func (s *VerificationService) SetEmailVerified(ctx context.Context, email string) error {

	// update the user table isVerified to true
	if err := s.repo.MarkEmailVerified(ctx, email); err != nil {
		return err
	}

	return nil

}
