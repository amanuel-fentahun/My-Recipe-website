package actions

import (
	"go-functions/internal/response"

	"github.com/gin-gonic/gin"
)

type ForgotPassowrdPayload struct {
	Input struct {
		Arg1 struct {
			Email string `json:"email"`
		} `json:"arg1"`
	} `json:"input"`
}

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

func LoginHandler(c *gin.Context) {

}

func SignUpHandler(c *gin.Context) {

}

func ForgotPasswordHandler(c *gin.Context) {

	var payload ForgotPassowrdPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(response.NewValidationError("Invalid Email Input", err))
		return
	}

	err := authService.InitiatePasswordReset(c.Request.Context(), payload.Input.Arg1.Email)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendOk(c, gin.H{
		"message": "Password reset code has been sent. Please check your email.",
	})
}

func PasswordResetHandler(c *gin.Context) {

	var payload ResetPasswordPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(response.NewValidationError("Invalid Password reset input", err))
		return
	}

}
