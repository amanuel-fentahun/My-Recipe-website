package auth

import (
	"errors"
	"time"
)

func Verify(actualCode string, providedCode string, expiry time.Time) (bool, error) {

	// check the expiration

	if expiry.UTC().Before(time.Now().UTC()) {
		return false, errors.New("verification code has expired")
	}

	// check the code
	if actualCode != providedCode {
		return false, errors.New("invalid verification code")
	}

	return true, nil

}
