package utils

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
)

func Credentails() (*cloudinary.Cloudinary, context.Context) {
	cld, _ := cloudinary.New()
	cld.Config.URL.Secure = true
	cxt := context.Background()

	return cld, cxt
}

func GenerateUploadSignature(params map[string]string, apiSecret string) string {

	//collect keys in alphabetical order
	keys := make([]string, 0, len(params))

	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	//Build "key=value" strings
	var sb strings.Builder

	for i, k := range keys {
		if params[k] != "" {
			sb.WriteString(fmt.Sprintf("%s=%s", k, params[k]))

			if i < len(keys)-1 {
				sb.WriteString("&")
			}
		}
	}

	//Append api_secret at the end
	toSign := sb.String() + apiSecret

	//Hash the resulting string using SHA-1
	h := sha1.New()
	h.Write([]byte(toSign))
	return hex.EncodeToString(h.Sum(nil))
}
