package routes

import (
	cloudinaryhandler "go-functions/Handlers/Cloudinary_handler"
	verifyemail "go-functions/Handlers/Email_verification_handler"
	hasuraactionhandler "go-functions/Handlers/Hasura_action_handler"
	passresethandler "go-functions/Handlers/Password-reset_handler"
	securitymiddleware "go-functions/Middlewares/security"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {

	router.Use(securitymiddleware.ValidateIncomingRequest)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.POST("/login", hasuraactionhandler.LoginHandler)
	router.GET("/upload_signature", cloudinaryhandler.CloudinarySignatureHandler)
	router.POST("/signup", hasuraactionhandler.SignUpHandler)

	// both API endpoints for password reset verification and email verification
	{
		verify := router.Group("")
		verify.Use(securitymiddleware.ValidateVerificationData)
		verify.POST("/verify_email", verifyemail.VerifyEmailHandler)
		verify.POST("/verify_pass_reset_code", passresethandler.VerifyPassResetCode)
	}

	router.POST("/forgot_password", passresethandler.GeneratePassResetCode)
	router.POST("/password_reset", passresethandler.ResetPasswordHandler)
}
