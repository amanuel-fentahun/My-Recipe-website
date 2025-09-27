package passwordresethandler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go-functions/config"
	"go-functions/internal/mail"
	"go-functions/internal/random"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
)

type ForgotPassowrdPayload struct {
	Input struct {
		Arg1 struct {
			Email string `json:"email"`
		} `json:"arg1"`
	} `json:"input"`
}

func GeneratePassResetCode(c *gin.Context) {

	var payload ForgotPassowrdPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	email := payload.Input.Arg1.Email

	if len(email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email input required!"})
		return
	}

	var query struct {
		Users []struct {
			Id uuid.UUID `graphql:"id"`
		} `graphql:"Users(where: {email: {_eq: $email}})"`
	}

	vars := map[string]interface{}{
		"email": graphql.String(email),
	}

	if err := config.NewGraphqlClient().Query(context.Background(), &query, vars); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server error. please try again later!"})
		return
	}

	if len(query.Users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
		return
	}

	code := random.GenerateRandomString(6)

	// store the verification code in the database for later verification
	var VerificationDataMutation struct {
		InsertVerificationDataOne struct {
			Email string `graphql:"email"`
		} `graphql:"insert_VerificationData_one(object: {email: $email, code: $code, type: $type})"`
	}

	vars2 := map[string]interface{}{
		"email": graphql.String(email),
		"code":  graphql.String(code),
		"type":  graphql.String("email_verification"),
	}

	if err := config.NewGraphqlClient().Mutate(context.Background(), &VerificationDataMutation, vars2); err != nil {
		c.JSON(500, gin.H{"error": "Unable to send password reset code. Please try again later"})
		return
	}

	//Send A verification Email

	subject := "Password reset code"
	body := "<p> You need to insert this code in order to keep resetting your password.</p>"
	body += "<p>Your password reset code:</p>"
	body += fmt.Sprintf("<h3>%s</h3>", code)
	body += "<p>If you did not request this password reset code, please ignore this email.</p>"
	body += "<p>Thanks,<br/>The Tafach Kitchen Team</p>"

	if err := mail.SendEmail(email, subject, body); err != nil {
		c.JSON(500, gin.H{
			"error": "Unable to send password reset code. Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Please check your mail for password reset code"})
}
