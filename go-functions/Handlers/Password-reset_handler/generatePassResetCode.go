package passwordresethandler

import (
	"context"
	"net/http"

	"go-functions/config"

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
	}

	email := payload.Input.Arg1.Email

	if len(email) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email input required!"})
	}

	var query struct {
		Users []struct {
			Id uuid.UUID `json:"id"`
		} `json:"Users(where: {email: {_eq: $email}})"`
	}

	vars := map[string]interface{}{
		"email": graphql.String(email),
	}

	if err := config.NewGraphqlClient().Query(context.Background(), &query, vars); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get this user information! please try again!"})
	}

}
