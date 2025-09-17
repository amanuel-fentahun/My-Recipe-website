package emailverificationhandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
)

type VerificationData struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expireAt"`
	Type     string    `json:"type"`
}

func VerifyEmailHandler(c *gin.Context) {

	email := c.MustGet("email").(string)
	code := c.MustGet("code").(string)

	client := graphql.NewClient(os.Getenv("HASURA_GRAPHQL_ENDPOINT"), nil).
		WithRequestModifier(func(r *http.Request) {
			r.Header.Set("x-hasura-admin-secret", os.Getenv("HASURA_ADMIN_SECRET"))
		})

	var query struct {
		VerificationData `graphql:"VerificationData_by_pk(email: $email)"`
	}

	vars := map[string]interface{}{
		"email": email,
	}

	if err := client.Query(context.Background(), &query, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch verification data"})
		fmt.Println(err)
		return
	}

	if query.VerificationData.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No verification data found for this email"})
		return
	}

	valid, err := verify(query.VerificationData.Code, code, query.VerificationData.ExpireAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if valid {
		// update the user table isVerified to true
		c.JSON(http.StatusOK, gin.H{"message": "Your email is verifed successfully. Please login to your account to see your profile! Thank you."})
		return
	}
}

func verify(actualCode string, providedCode string, expiry time.Time) (bool, error) {

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
