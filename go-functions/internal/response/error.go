package response

import (
	"fmt"
	"net/http"
)

type BusinessCode string

const (
	CodeAuthFailed        BusinessCode = "AUTH_UNAUTHORIZED"
	CodePermissionDenied  BusinessCode = "AUTH_FORBIDDEN"
	CodeNotificationError BusinessCode = "NOTIFICATION_FAILED"
	CodeCloudinaryError   BusinessCode = "CLOUDINARY_SIGN_FAILED"
	CodeInvalidInput      BusinessCode = "BAD_REQUEST"
	CodeInternalError     BusinessCode = "INTERNAL_SERVER_ERROR"
)

type AppError struct {
	HTTPStatus int          `json:"-"`
	Code       BusinessCode `json:"code"`
	Message    string       `json:"message"`
	RawError   error        `json:"-"`
}

func (e *AppError) Error() string {
	if e.RawError != nil {
		return fmt.Sprintf("[%s] %s: $v", e.Code, e.Message, e.RawError)
	}

	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.RawError
}

func NewUnauthorizedError(msg string, err error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       CodeAuthFailed,
		Message:    msg,
		RawError:   err,
	}
}

func NewForbiddenError(msg string, err error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusForbidden,
		Code:       CodePermissionDenied,
		Message:    msg,
		RawError:   err,
	}
}

func NewNotificationError(msg string, err error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       CodeNotificationError,
		Message:    msg,
		RawError:   err,
	}
}

func NewCloudinaryError(msg string, err error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		Code:       CodeCloudinaryError,
		Message:    msg,
		RawError:   err,
	}
}

func NewValidationError(msg string, err error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		Code:       CodeInvalidInput,
		Message:    msg,
		RawError:   err,
	}
}
