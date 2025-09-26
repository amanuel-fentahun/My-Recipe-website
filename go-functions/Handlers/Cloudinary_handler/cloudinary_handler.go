package cloudinaryhandler

import (
	"fmt"
	"os"
	"time"

	"go-functions/internal/cloud"

	"github.com/gin-gonic/gin"
)

func CloudinarySignatureHandler(c *gin.Context) {
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	params := map[string]string{
		"timestamp": timestamp,
		"eager":     "w_400,h_300,c_pad|w_260,h_200,c_crop",
		"public_id": "sample_image",
	}

	signature := cloud.GenerateUploadSignature(params, apiSecret)

	fmt.Println("Signature: ", signature)
	fmt.Println("Timestamp: ", timestamp)
}
