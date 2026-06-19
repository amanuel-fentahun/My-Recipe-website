package actions

import "github.com/gin-gonic/gin"

func VerifyEmailHandler(c *gin.Context)
{
	expired_at := c.MustGet("expired_at").((time.Time))
	old_code := c.MustGet("old_code").(string)
	incoming_code := c.MustGet("incoming_code").(string)
	email := c.MustGet("email").(string)
}