package passwordresethandler

import (
	"context"
	"go-functions/config"
	"go-functions/internal/auth"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hasura/go-graphql-client"
)

type ResetPasswordPayload struct {
	Input struct {
		Inputs struct {
			Email              string `json:"email"`
			NewPassword        string `json:"newPassword"`
			ConfirmNewPassowrd string `json:"confirmNewPassowrd"`
			SecretCode         string `json:"secretCode"`
		} `json:"inputs"`
	} `json:"input"`
}

type VerificationData struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expireAt"`
	Type     string    `json:"type"`
}

func ResetPasswordHandler(c *gin.Context) {

	var payload ResetPasswordPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Println("unable to decode password reset request json")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	inputs := payload.Input.Inputs

	var query struct {
		VerificationData `graphql:"VerificationData_by_pk(email: $email)"`
	}

	vars := map[string]interface{}{
		"email": inputs.Email,
	}

	if err := config.NewGraphqlClient().Query(context.Background(), &query, vars); err != nil {
		log.Println("unable to retrieve the verification data from db")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reset password. Please try again later"})
		return
	}

	if query.VerificationData.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to reset password. Please try again later"})
		return
	}

	if query.VerificationData.Code != inputs.SecretCode {
		log.Println("verification code did not match even after verifying PassReset")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reset password. Please try again later"})
		return
	}

	if inputs.ConfirmNewPassowrd != inputs.NewPassword {
		log.Println("password and password re-enter did not match")
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "password and confirm password are not the same"})
		return
	}

	hashed_password, err := auth.HashPassword(inputs.NewPassword)

	if err != nil {
		log.Println("Failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reset password. Please try again later"})
	}

	// update the users password
	var Mutation struct {
		UpdateUsers struct {
			Affected_rows int `graphql:"affected_rows"`
		} `graphql:"update_Users(where: {email: {_eq: $email}}, _set:{password: $password})"`
	}

	vars2 := map[string]interface{}{
		"email":    graphql.String(inputs.Email),
		"password": graphql.String(hashed_password),
	}

	if err := config.NewGraphqlClient().Mutate(context.Background(), &Mutation, vars2); err != nil {
		log.Println("unable to update user password in db")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reset password. Please try again later"})
		return
	}

	if Mutation.UpdateUsers.Affected_rows == 0 {
		log.Println("no affected row during the mutation of updating the password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to reset password. Please try again later"})
		return
	}

	// delete teh verification data
	var Mutation2 struct {
		DeleteVerificationData struct {
			Email string `graphql:"email"`
		} `graphql:"delete_VerificationData_by_pk(email: $email)"`
	}

	vars3 := map[string]interface{}{
		"email": graphql.String(inputs.Email),
	}

	if err := config.NewGraphqlClient().Mutate(context.Background(), &Mutation2, vars3); err != nil {
		log.Println("Unable to delete the verification data from db")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong. Please try again later!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password Reset Successfully"})
}
