package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VerificationData struct {
	Input struct {
		Inputs struct {
			Email string `json:"email"`
			Code  string `json:"code"`
		} `json:"verCode"`
	} `json:"input"`
}

func ValidateVerificationData(c *gin.Context) {

	var payload VerificationData

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	input := payload.Input.Inputs

	if len(input.Email) == 0 || len(input.Code) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "Email and Code are required"})
		return
	}

	c.Set("email", input.Email)
	c.Set("code", input.Code)

	c.Next()
}
