package services

import (
	"go-functions/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.HasuraRepository
	jwtSecret []byte
}

type UserProfile struct {
	fullName  string
	email     string
	userId    uuid.UUID
	avaterURL string
	roles     []string
}

func NewAuthService(repo *repository.HasuraRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret)}
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
