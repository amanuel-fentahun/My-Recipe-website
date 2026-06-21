package actions

import (
	"go-functions/internal/response"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {

}

func SignUpHandler(c *gin.Context) {

}

func ForgotPasswordHandler(c *gin.Context) {

}

type ForgotPassowrdPayload struct {
	Input struct {
		Arg1 struct {
			Email string `json:"email"`
		} `json:"arg1"`
	} `json:"input"`
}

func PasswordResetHandler(c *gin.Context) {

	var payload ForgotPassowrdPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(response.NewValidationError("Invalid Email Input", err))
		return
	}
}
