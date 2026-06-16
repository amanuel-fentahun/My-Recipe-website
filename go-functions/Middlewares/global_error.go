package middlewares

import (
	"errors"
	"log"
	"net/http"

	"go-functions/internal/response"
	myError "go-functions/internal/response"

	"github.com/gin-gonic/gin"
)

type HasuraErrorPayload struct {
	Message string               `json:"message"`
	Code    myError.BusinessCode `json:"code"`
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {

			lastErr := c.Errors.Last().Err

			var appErr *myError.AppError
			if errors.As(lastErr, &appErr) {

				if appErr.RawError != nil {
					log.Printf("[Error] Code: %s | Message: %s | Internal: %v", appErr.Code, appErr.Message, appErr.RawError)
				}

				c.JSON(appErr.HTTPStatus, HasuraErrorPayload{
					Message: appErr.Message,
					Code:    appErr.Code,
				})

			} else {

				log.Printf("[UNHANDLED ERROR] %v", lastErr)
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
	return gin.CustomRecovery(func(c *gin.Context, err any) {

		log.Printf("[PANIC RECOVERED] Critical System Crash: %v", err)

		c.JSON(http.StatusInternalServerError, response.StandardResponse{
			Success: false,
			Message: "A critical internal server error occurred",
		})

		c.Abort()
	})
}
