package hasuraactionhandler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-functions/config"
	"go-functions/internal/auth"
	"go-functions/internal/mail"
	"go-functions/internal/random"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
)

type SignupInput struct {
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
}

type VerificationData struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expireAt"`
	Type     string    `json:"type"`
}

type SignupActionPayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            struct {
		Inputs SignupInput `json:"inputs"`
	} `json:"input"`
}

func SignUpHandler(c *gin.Context) {

	var actionPayload SignupActionPayload

	if err := c.ShouldBindJSON(&actionPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	input := actionPayload.Input.Inputs

	// GraphQL query for select one user by email
	var query struct {
		UsersAggrigation struct {
			Nodes []struct {
				ID uuid.UUID `graphql:"id"`
			} `graphql:"nodes"`
		} `graphql:"Users_aggregate(where: {email: {_eq: $email}})"`
	}

	vars := map[string]interface{}{
		"email": graphql.String(input.Email),
	}

	// try to fetch user with the given email
	if err := config.NewGraphqlClient().Query(context.Background(), &query, vars); err != nil {
		c.JSON(500, gin.H{"error": "Failed to check existing user"})
		return
	}

	// check if the user already exists
	if len(query.UsersAggrigation.Nodes) > 0 {
		c.JSON(400, gin.H{"error": "User with this email already exists"})
		return
	}

	// hash the user password
	hashed_password, err := auth.HashPassword(input.Password)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// Graphql mutation for signup
	var SignupMutation struct {
		InsertUsersOne struct {
			ID uuid.UUID `graphql:"id"`
		} `graphql:"insert_Users_one(object: {name: $name, email: $email, password: $password, avater_url: $avater_url, isVerified: false})"`
	}

	vars2 := map[string]interface{}{
		"name":       graphql.String(input.FullName),
		"email":      graphql.String(input.Email),
		"password":   graphql.String(hashed_password),
		"avater_url": graphql.String(input.AvatarURL),
	}

	// try to create a new user
	if err := config.NewGraphqlClient().Mutate(context.Background(), &SignupMutation, vars2); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	code := random.GenerateRandomString(6)

	// store the verification code in the database for later verification
	var VerificationDataMutation struct {
		InsertVerificationDataOne struct {
			Email string `graphql:"email"`
		} `graphql:"insert_VerificationData_one(object: {email: $email, code: $code, type: $type})"`
	}

	vars3 := map[string]interface{}{
		"email": graphql.String(input.Email),
		"code":  graphql.String(code),
		"type":  graphql.String("email_verification"),
	}

	if err := config.NewGraphqlClient().Mutate(context.Background(), &VerificationDataMutation, vars3); err != nil {
		c.JSON(500, gin.H{"error": "Failed to store verification data"})
		return
	}

	//Send A verification Email

	subject := "Verify your Email"
	body := "<p> You need to verfiy your email address to continue using your <strong>Tafach Kitchen</strong> account.</p>"
	body += "<p>Enter the following code to verify your email address:</p>"
	body += fmt.Sprintf("<h3>%s</h3>", code)
	body += "<p>If you did not create this account, please ignore this email.</p>"
	body += "<p>Thanks,<br/>The Tafach Kitchen Team</p>"

	if err := mail.SendEmail(input.Email, subject, body); err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to send verification email",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User created successfully, Please verify your email",
	})

}
