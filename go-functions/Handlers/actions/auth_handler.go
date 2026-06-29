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

type ResendCodePayload struct {
	Input struct {
		Arg1 struct {
			Email      string `json:"email"`
			ActionType string `json:"actionType"`
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

	inputs := payload.Input.Inputs

	if err := authService.CompletePasswordReset(
		c.Request.Context(),
		inputs.Email,
		inputs.SecretCode,
		inputs.NewPassword,
		inputs.ConfirmNewPassowrd,
	); err != nil {
		_ = c.Error(err)
		return
	}

	response.SendOk(c, gin.H{
		"message": "Your password is resetted seccussfully.",
	})
}

func ResendVerificationCodeHandler(c *gin.Context) {
	var payload ResendCodePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(response.NewValidationError("Invalid email input.", err))
		return
	}

	err := authService.InitiateVerificationSend(c.Request.Context(), payload.Input.Arg1.Email, payload.Input.Arg1.ActionType)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendOk(c, gin.H{
		"message": "A new verification code has been send to your email address.",
		"status":  "SUCCESS",
	})
}
