package security

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ValidateIncomingRequest(c *gin.Context) {
	incomingEventSecret := c.GetHeader("x-hasura-event-secret")

	eventSecret := os.Getenv("EVENT_SECRET")

	if incomingEventSecret != eventSecret {
		log.Println("Unauthorized request: invalid or missing event secret")
		c.Abort()
		c.JSON(http.StatusForbidden, gin.H{"error": "untrust origin not supported"})
		return
	}

	c.Next()
}
