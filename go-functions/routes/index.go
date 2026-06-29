package routes

import (
	"errors"
	// cloudinaryhandler "go-functions/Handlers/Cloudinary_handler"
	// verifyemail "go-functions/Handlers/Email_verification_handler"
	// hasuraactionhandler "go-functions/Handlers/Hasura_action_handler"
	//passresethandler "go-functions/Handlers/Password-reset_handler"
	actionHandler "go-functions/Handlers/actions"
	middlewares "go-functions/Middlewares"
	"go-functions/internal/response"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {

	router.Use(middlewares.CustomRecovery())
	router.Use(middlewares.GlobalErrorHandler())

	// Global 404 Route catcher
	router.NoRoute(func(c *gin.Context) {
		err := errors.New("the requested endpoint could not be found")
		_ = c.Error(response.NewValidationError("Route not found", err))
	})

	// public Endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Hasura Protected Group (Requires the Hasura Event Secret Key)
	hasuraProtected := router.Group("")
	hasuraProtected.Use(middlewares.ValidateIncomingRequest())
	{
		// hasuraProtected.POST("/login", hasuraactionhandler.LoginHandler)
		// hasuraProtected.GET("/upload_signature", cloudinaryhandler.CloudinarySignatureHandler)
		//hasuraProtected.POST("/signup", hasuraactionhandler.SignUpHandler)
		hasuraProtected.POST("/forgot_password", actionHandler.ForgotPasswordHandler)
		hasuraProtected.POST("/password_reset", actionHandler.PasswordResetHandler)
		hasuraProtected.POST("/resend_code", actionHandler.ResendVerificationCodeHandler)

		// (Requires BOTH the Hasura Secret AND the Verification Data checked)
		verificationGroup := hasuraProtected.Group("")
		verificationGroup.Use(middlewares.ValidateVerificationData())
		{
			verificationGroup.POST("/verify_email", actionHandler.VerifyEmailHandler)
			verificationGroup.POST("/verify_resetCode", actionHandler.VerifyResetCode)
		}
	}

}
