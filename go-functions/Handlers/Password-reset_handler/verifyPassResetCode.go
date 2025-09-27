package passwordresethandler

import (
	"go-functions/internal/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func VerifyPassResetCode(c *gin.Context) {

	expired_at := c.MustGet("expired_at").((time.Time))
	old_code := c.MustGet("old_code").(string)
	incoming_code := c.MustGet("incoming_code").(string)

	valid, err := auth.Verify(old_code, incoming_code, expired_at)

	if !valid {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Password Reset code verification succeeded",
		"data":    old_code})

}
