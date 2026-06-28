package services

import (
	"context"
	utils "go-functions/Utils"
	"go-functions/internal/mail"
	"go-functions/internal/repository"
	"go-functions/internal/response"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.HasuraRepository
	jwtSecret []byte
	codeTTL   time.Duration
}

type UserProfile struct {
	fullName  string
	email     string
	userId    uuid.UUID
	avaterURL string
	roles     []string
}

func NewAuthService(repo *repository.HasuraRepository, jwtSecret string, codeTTLStr string) *AuthService {
	duration := 15 * time.Minute

	if codeTTLStr != "" {
		parsed, err := time.ParseDuration(codeTTLStr)
		if err != nil {
			log.Printf("[CONFIG WARNING] Invalid VERIFICATION_CODE_TTL value '%s'. Falling back to 15m. Error: %v", codeTTLStr, err)
		} else {
			duration = parsed
		}
	}
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		codeTTL:   duration,
	}

}

func (s *AuthService) CreateToken(userFromDB UserProfile) (string, error) {

	defaultRoles := "user"
	if len(userFromDB.roles) > 0 {
		defaultRoles = userFromDB.roles[0]
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"fullName":  userFromDB.fullName,
			"email":     userFromDB.email,
			"userId":    userFromDB.userId,
			"avaterURL": userFromDB.avaterURL,
			"https://hasura.io/jwt/claims": map[string]interface{}{
				"x-hasura-allowed-roles": userFromDB.roles,
				"x-hasura-default-role":  defaultRoles,
				"x-hasura-user-id":       userFromDB.userId.String(),
			},
			"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		})

	return jwtToken.SignedString(s.jwtSecret)
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(plainPassword, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
	return err == nil
}

func (s *AuthService) InitiatePasswordReset(ctx context.Context, email string) error {

	if err := s.InitiateVerificationSend(ctx, email, string(repository.ActionPasswordReset)); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) CompletePasswordReset(ctx context.Context, email, secretCode, newPassword, confirmpassword string) error {

	if newPassword != confirmpassword {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "Password do not match one another.",
		}
	}

	if len(newPassword) < 8 {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "Your new password is too short. It must consist of 8 or more characters.",
		}
	}

	if !utils.IsValidEmail(email) {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "Please provide a valid email address.",
		}
	}

	verification, err := s.repo.FetchVerificationDataByEmail(ctx, email)
	if err != nil {
		return err
	}

	if verification.Code != secretCode {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "The password reset is failed. Please try again.",
		}
	}

	if time.Now().After(verification.ExpireAt) {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "The password reset token has expired.",
		}
	}

	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return &response.AppError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       response.CodeInternalError,
			Message:    "Internal server error.",
			RawError:   err,
		}
	}

	if err := s.repo.UpdateUserPassword(ctx, email, hashedPassword); err != nil {
		return err
	}

	if err := s.repo.ArchiveAndPurgeVerificationRow(ctx, email, secretCode, string(repository.ActionPasswordReset), "SUCCESS"); err != nil {
		log.Printf("[WARNING] Audit log processing sequence encountered an interruption: %v", err)
	}

	return nil
}

func (s *AuthService) InitiateVerificationSend(ctx context.Context, email, actionType string) error {

	if !utils.IsValidEmail(email) {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "Please provide a valid email address.",
		}
	}

	userExists, err := s.repo.CheckIfUserExists(ctx, email)

	if err != nil {
		return err
	}

	if !userExists {
		log.Printf("[SECURITY] Verification Code send requested for non-existent email: %s", email)
		return nil
	}

	status, currentData, err := s.repo.CheckVerificationState(ctx, email)
	if err != nil {
		return err
	}

	if status == repository.StatusActiveCode {
		return &response.AppError{
			HTTPStatus: http.StatusTooManyRequests,
			Code:       response.CodeRateLimitExceeded,
			Message:    "Your code is already active. Please wait a moment before asking for another.",
		}
	}

	if status != repository.StatusNoRowExists {
		if err := s.repo.ArchiveAndPurgeVerificationRow(ctx, email, currentData.Code, actionType, "EXPIRED"); err != nil {
			log.Printf("[WARNING] Audit log processing sequence encountered an interruption: %v", err)
		}
	}

	newCode := utils.GenerateRandomString(6)

	err = s.repo.InsertVerificationRow(ctx, email, newCode, s.codeTTL, actionType)
	if err != nil {
		return err
	}

	var subject, body string
	if actionType == string(repository.ActionPasswordReset) {
		subject = utils.SubjectPasswordReset
		body = utils.GetPasswordResetTemplate(newCode)
	} else {
		subject = utils.SubjectEmailVerification
		body = utils.GetEmailVerificationTemplate(newCode)
	}

	err = mail.SendEmail(email, subject, body)
	if err != nil {
		return response.NewSMTPMailError("we could not send your code. Please try again later.", err)
	}

	return nil

}
