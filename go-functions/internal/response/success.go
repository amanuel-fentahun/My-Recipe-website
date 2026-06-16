package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SendOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Data:    data,
	})
}

func SendCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SendDeleted(c *gin.Context, message string) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
	})
}
