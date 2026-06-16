package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)

	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
