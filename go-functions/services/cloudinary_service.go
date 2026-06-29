package services

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"go-functions/internal/response"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CoudinaryService struct {
	apiKey    string
	apiSecret string
}

func NewCloudinaryService() *CoudinaryService {
	return &CoudinaryService{
		apiKey:    os.Getenv("CLOUDINARY_API_KEY"),
		apiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}
}

func (s *CoudinaryService) GenerateUploadSignature(folder string) (gin.H, error) {
	if s.apiKey == "" || s.apiSecret == "" {
		err := errors.New("missing production cloudinary environment credentials configuration keys")
		return nil, &response.AppError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       response.CodeCloudinaryError,
			Message:    "Cloud media storage integration layer is temporarily offline.",
			RawError:   err,
		}
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	params := map[string]string{
		"folder":    folder,
		"timestamp": timestamp,
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		if params[k] != "" {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, params[k]))
		}
	}

	toSign := strings.Join(pairs, "&") + s.apiSecret

	h := sha1.New()
	h.Write([]byte(toSign))
	signature := hex.EncodeToString(h.Sum(nil))

	return gin.H{
		"signature": signature,
		"timestamp": timestamp,
		"apiKey":    s.apiKey,
	}, nil
}
