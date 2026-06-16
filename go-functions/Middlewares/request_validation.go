package middlewares

import (
	"errors"
	"go-functions/internal/response"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type VerificationDataPayload struct {
	Input struct {
		Inputs struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		} `json:"verCode`
	} `json:"input"`
}

type VerificationData struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expireAt"`
	Type     string    `json:"type"`
}

func ValidateIncomingRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		incomingEventSecret := c.GetHeader("x-hasura-event-secret")
		eventSecret := os.Getenv("EVENT_SECRET")

		if incomingEventSecret == "" || incomingEventSecret != eventSecret {
			log.Println("[SECURITY] Unauthorized request attempt: invalid or missing event secret")

			err := errors.New("even secrect mismatch or header missing")

			appErr := response.NewForbiddenError("untrust origin not supported", err)

			_ = c.Error(appErr)
			c.Abort()
			return
		}

		c.Next()
	}
}

func ValidateVerificationData() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload VerificationDataPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			appErr := response.NewValidationError("Invalid request payload format structure", err)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

		input := payload.Input.Inputs

		if len(input.Email) == 0 || len(input.Code) == 0 {
			err := errors.New("missing email or verification code parameter fields")
			appErr := response.NewValidationError("Email and Code fields are required", err)
			_ = c.Error(appErr)
			c.Abort()
			return
		}

	}
}
