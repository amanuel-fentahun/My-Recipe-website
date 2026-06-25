package middlewares

import (
	"errors"
	utils "go-functions/Utils"
	"go-functions/internal/repository"
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
		} `json:"verCode"`
	} `json:"input"`
}

type VerificationData struct {
	Email    string    `json:"email"`
	Code     string    `json:"code"`
	ExpireAt time.Time `json:"expireAt"`
	Type     string    `json:"type"`
}

var HasuraRepo = repository.NewHasuraRepository()

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
			_ = c.Error(response.NewValidationError("Invalid request payload format structure", err))
			c.Abort()
			return
		}

		input := payload.Input.Inputs

		if !utils.IsValidEmail(input.Email) {
			err := errors.New("Invalid email address")
			_ = c.Error(response.NewValidationError("Email field is required", err))
			c.Abort()
			return
		}

		if len(input.Code) == 0 {
			err := errors.New("missing verification code parameter field")
			_ = c.Error(response.NewValidationError("Code field are required", err))
			c.Abort()
			return
		}

		data, err := HasuraRepo.FetchVerificationDataByEmail(c.Request.Context(), input.Email)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}

		c.Set("old_code", data.Code)
		c.Set("expired_at", data.ExpireAt)
		c.Set("incoming_code", input.Code)
		c.Set("email", data.Email)

		c.Next()

	}
}
