package routes

import (
	cloudinaryhandler "go-functions/Handlers/Cloudinary_handler"
	hasuraactionhandler "go-functions/Handlers/Hasura_action_handler"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.POST("/login", hasuraactionhandler.LoginHandler)
	router.GET("/get-upload-signature", cloudinaryhandler.CloudinarySignatureHandler)
	router.POST("/signup", hasuraactionhandler.SignUpHandler)
}
