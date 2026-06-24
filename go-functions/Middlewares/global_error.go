package middlewares

import (
	"errors"
	"log"
	"net/http"
	"os"

	myError "go-functions/internal/response"

	"github.com/gin-gonic/gin"
)

type HasuraErrorPayload struct {
	Message string               `json:"message"`
	Code    myError.BusinessCode `json:"code"`
}

func GlobalErrorHandler() gin.HandlerFunc {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {

			lastErr := c.Errors.Last().Err

			var appErr *myError.AppError
			if errors.As(lastErr, &appErr) {

				if appErr.RawError != nil {
					logger.Printf("[BACKEND EXCEPTION LOG] BusinessCode: %s | SafeMessage: %s | RAW SYSTEM ERROR: %v", appErr.Code, appErr.Message, appErr.RawError)
				}

				c.JSON(appErr.HTTPStatus, HasuraErrorPayload{
					Message: appErr.Message,
					Code:    appErr.Code,
				})

			} else {

				logger.Printf("[UNHANDLED ERROR] %v", lastErr)
				c.JSON(http.StatusInternalServerError, HasuraErrorPayload{
					Message: "An unexpected internal server error occurred",
					Code:    myError.CodeInternalError,
				})
			}

			c.Abort()
		}
	}
}

func CustomRecovery() gin.HandlerFunc {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	return gin.CustomRecovery(func(c *gin.Context, err any) {

		logger.Printf("[PANIC RECOVERED] Critical System Crash: %v", err)

		c.JSON(http.StatusInternalServerError, myError.StandardResponse{
			Success: false,
			Message: "A critical internal server error occurred",
		})

		c.Abort()
	})
}
