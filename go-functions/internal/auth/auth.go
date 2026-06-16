package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserProfile struct {
	fullName  string
	email     string
	userId    uuid.UUID
	avaterURL string
	roles     []string
}

func CreateToken(userFromDB UserProfile, jwtSecret []byte) (string, error) {

	var defaultRoles string
	if len(userFromDB.roles) > 0 {
		defaultRoles = userFromDB.roles[0]
	} else {
		defaultRoles = "user"
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

	return jwtToken.SignedString(jwtSecret)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Verify(actualCode string, providedCode string, expiry time.Time) (bool, error) {

	// check the expiration

	if expiry.UTC().Before(time.Now().UTC()) {
		return false, errors.New("verification code has expired")
	}

	// check the code
	if actualCode != providedCode {
		return false, errors.New("invalid verification code")
	}

	return true, nil

}
