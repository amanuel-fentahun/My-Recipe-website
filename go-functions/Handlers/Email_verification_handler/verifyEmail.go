package emailverificationhandler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go-functions/config"

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

	var query struct {
		VerificationData `graphql:"VerificationData_by_pk(email: $email)"`
	}

	vars := map[string]interface{}{
		"email": email,
	}

	if err := config.NewGraphqlClient().Query(context.Background(), &query, vars); err != nil {
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
		var Mutation struct {
			UpdateUsers struct {
				Affected_rows int `graphql:"affected_rows"`
			} `graphql:"update_Users(where: {email: {_eq: $email}}, _set: {isVerified: true})"`
		}

		// delete the verificationData after succussfully verify the email
		var Mutation2 struct {
			DeleteVerificationData struct {
				Email string `graphql:"email"`
			} `graphql:"delete_VerificationData_by_pk(email: $email)"`
		}

		vars2 := map[string]interface{}{
			"email": graphql.String(email),
		}

		if err := config.NewGraphqlClient().Mutate(context.Background(), &Mutation, vars2); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong. Please try again later!"})
			return
		}

		// if isVerified updated successfully delete the verificationData
		if Mutation.UpdateUsers.Affected_rows == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to verify your email. Please try again later!"})
			return
		}

		if err := config.NewGraphqlClient().Mutate(context.Background(), &Mutation2, vars2); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong. Please try again later!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Your email is verifed successfully. Please login to your account to access your profile! Thank you."})
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
