package actions

import (
	"go-functions/internal/response"
	"time"

	"github.com/gin-gonic/gin"
)

func VerifyEmailHandler(c *gin.Context) {
	expired_at := c.MustGet("expired_at").((time.Time))
	old_code := c.MustGet("old_code").(string)
	incoming_code := c.MustGet("incoming_code").(string)
	email := c.MustGet("email").(string)

	if err := verifyService.VerifyEmail(c.Request.Context(), email, old_code, incoming_code, expired_at); err != nil {
		_ = c.Error(err)
		return
	}

	if err := verifyService.SetEmailVerified(c.Request.Context(), old_code, email); err != nil {
		_ = c.Error(err)
		return
	}

	response.SendOk(c, gin.H{
		"status":  "verified",
		"message": "Your email address has been successfully verified.",
	})
}

func VerifyResetCode(c *gin.Context) {
	expired_at := c.MustGet("expired_at").((time.Time))
	old_code := c.MustGet("old_code").(string)
	incoming_code := c.MustGet("incoming_code").(string)
	email := c.MustGet("email").(string)

	if err := verifyService.VerifyEmail(c.Request.Context(), email, old_code, incoming_code, expired_at); err != nil {
		_ = c.Error(err)
		return
	}

	response.SendOk(c, gin.H{
		"status":  "Valid",
		"message": "Your password reset code is valid. Continue resetting.",
	})
}
