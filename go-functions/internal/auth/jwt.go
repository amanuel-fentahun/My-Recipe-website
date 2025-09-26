package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateToken(name string, email string, userId uuid.UUID, avaterURL string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"fullName":  name,
			"email":     email,
			"userId":    userId,
			"avaterURL": avaterURL,
			"https://hasura.io/jwt/claims": map[string]interface{}{
				"x-hasura-allowed-roles": []string{"user"},
				"x-hasura-default-role":  "user",
				"x-hasura-user-id":       userId.String(),
			},
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
