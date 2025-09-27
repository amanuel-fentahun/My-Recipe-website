package middlewares

import (
	"context"
	"go-functions/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type VerificationDataPayload struct {
	Input struct {
		Inputs struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		} `json:"verCode"`
	} `json:"input"`
}

type VerificationData struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expireAt"`
	Type     string    `json:"type"`
}

func ValidateVerificationData(c *gin.Context) {

	var payload VerificationDataPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	input := payload.Input.Inputs

	if len(input.Email) == 0 || len(input.Code) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Email and Code are required"})
		return
	}

	var query struct {
		VerificationData `graphql:"VerificationData_by_pk(email: $email)"`
	}

	vars := map[string]interface{}{
		"email": input.Email,
	}

	if err := config.NewGraphqlClient().Query(context.Background(), &query, vars); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again later"})
		return
	}

	if query.VerificationData.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "ErrUserNotFound"})
		return
	}

	c.Set("old_code", query.VerificationData.Code)
	c.Set("expired_at", query.VerificationData.ExpireAt)
	c.Set("incoming_code", input.Code)
	c.Set("email", input.Email)

	c.Next()
}
