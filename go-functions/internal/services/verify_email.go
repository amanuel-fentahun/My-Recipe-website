package services

import (
	"context"
	"net/http"
	"time"

	"go-functions/config"
	"go-functions/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
)

func VerifyEmailHandler(c *gin.Context) {

	expired_at := c.MustGet("expired_at").((time.Time))
	old_code := c.MustGet("old_code").(string)
	incoming_code := c.MustGet("incoming_code").(string)
	email := c.MustGet("email").(string)

	valid, err := auth.Verify(old_code, incoming_code, expired_at)

	if !valid {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
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
}
