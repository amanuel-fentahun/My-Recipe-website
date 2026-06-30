package actions

import (
	"go-functions/internal/response"

	"github.com/gin-gonic/gin"
)

type CloudinarySignaturePayload struct {
	Input struct {
		Arg1 struct {
			Folder string `json:"folder"`
		} `json:"arg1"`
	} `json:"input"`
}

func CloudinarySignatureHandler(c *gin.Context) {
	var payload CloudinarySignaturePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		_ = c.Error(response.NewValidationError("Invalid query signature format structure.", err))
		return
	}

	folder := payload.Input.Arg1.Folder
	if folder == "" {
		folder = "recipes"
	}

	signatureData, err := cloudinaryService.GenerateUploadSignature(folder)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendOk(c, gin.H{
		"data": signatureData,
	})
}
