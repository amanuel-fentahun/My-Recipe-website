package services

import (
	"context"
	"fmt"
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
		log.Printf("[SECURITY] Password reset requested for non-existent email: %s", email)
		return nil
	}

	status, _, err := s.repo.CheckVerificationState(ctx, email)
	if err != nil {
		return err
	}

	if status == "ACTIVE_CODE_WAIT" {
		return &response.AppError{
			HTTPStatus: http.StatusBadRequest,
			Code:       response.CodeInvalidInput,
			Message:    "A reset code is already active. Please wait a moment before asking for another.",
		}
	}

	newCode := utils.GenerateRandomString(6)
	err = s.repo.UpdateOrCreateVerificationRow(ctx, email, newCode, s.codeTTL, "password_reset")
	if err != nil {
		return err
	}

	subject := "Password reset code"
	body := "<p> You need to insert this code in order to keep resetting your password.</p>"
	body += "<p>Your password reset code:</p>"
	body += fmt.Sprintf("<h3>%s</h3>", newCode)
	body += "<p>If you did not request this password reset code, please ignore this email.</p>"
	body += "<p>Thanks,<br/>The Tafach Kitchen Team</p>"

	err = mail.SendEmail(email, subject, body)
	if err != nil {
		return response.NewSMTPMailError("we could not send your reset instructions. Please try again later.", err)
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
		if err := s.repo.ArchiveAndPurgeVerificationRow(ctx, email, secretCode, "password_reset", "EXPIRED"); err != nil {
			log.Printf("[WARNING] Audit log processing sequence encountered an interruption: %v", err)
		}

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

	if err := s.repo.ArchiveAndPurgeVerificationRow(ctx, email, secretCode, "password_reset", "SUCCESS"); err != nil {
		log.Printf("[WARNING] Audit log processing sequence encountered an interruption: %v", err)
	}

	return nil
}
