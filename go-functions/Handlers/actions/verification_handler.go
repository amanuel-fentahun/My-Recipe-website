package actions

import (
	"go-functions/internal/repository"
	"go-functions/internal/response"
	"go-functions/services"
	"time"

	"github.com/gin-gonic/gin"
)

var HasuraRepo = repository.NewHasuraRepository()
var verifyService = services.NewVerificationService(HasuraRepo)

func VerifyEmailHandler(c *gin.Context) {
	expired_at := c.MustGet("expired_at").((time.Time))
	old_code := c.MustGet("old_code").(string)
	incoming_code := c.MustGet("incoming_code").(string)
	email := c.MustGet("email").(string)

	if err := verifyService.VerifyEmail(c.Request.Context(), email, old_code, incoming_code, expired_at); err != nil {
		_ = c.Error(err)
		return
	}

	response.SendOk(c, gin.H{
		"status":  "verified",
		"message": "Your email address has been successfully verified.",
	})
}

func PasswordresetHandler(c *gin.Context) {

}
