package utils

import (
	"math/rand"
	"regexp"
	"strings"
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

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	cleanEmail := strings.TrimSpace(email)
	if cleanEmail == "" {
		return false
	}

	if len(cleanEmail) > 254 {
		return false
	}

	return emailRegex.MatchString(cleanEmail)
}
