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

// SendOk handles responses and flattens maps automatically to satisfy Hasura
func SendOk(c *gin.Context, data interface{}) {
	if data != nil {
		switch v := data.(type) {
		case gin.H:
			v["success"] = true
			c.JSON(http.StatusOK, v)
			return
		case map[string]interface{}:
			v["success"] = true
			c.JSON(http.StatusOK, v)
			return
		}
	}

	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Data:    data,
	})
}

// SendCreated handles successful creations and balances maps flatly
func SendCreated(c *gin.Context, message string, data interface{}) {
	if data != nil {
		switch v := data.(type) {
		case gin.H:
			v["success"] = true
			v["message"] = message
			c.JSON(http.StatusCreated, v)
			return
		case map[string]interface{}:
			v["success"] = true
			v["message"] = message
			c.JSON(http.StatusCreated, v)
			return
		}
	}

	c.JSON(http.StatusCreated, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendDeleted remains identical
func SendDeleted(c *gin.Context, message string) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
	})
}
