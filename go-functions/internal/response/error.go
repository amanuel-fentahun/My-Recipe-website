package response

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type BusinessCode string

const (
	CodeAuthFailed        BusinessCode = "AUTH_UNAUTHORIZED"
	CodePermissionDenied  BusinessCode = "AUTH_FORBIDDEN"
	CodeNotificationError BusinessCode = "NOTIFICATION_FAILED"
	CodeCloudinaryError   BusinessCode = "CLOUDINARY_SIGN_FAILED"
	CodeInvalidInput      BusinessCode = "BAD_REQUEST"
	CodeInternalError     BusinessCode = "INTERNAL_SERVER_ERROR"
	CodeSMTPError         BusinessCode = "SMTP_DELIVERY_FAILED"
	CodeRateLimitExceeded BusinessCode = "RATE_LIMIT_EXCEEDED"
	CodeDBError           BusinessCode = "DATABASE_ERROR"
)

type AppError struct {
	HTTPStatus int          `json:"-"`
	Code       BusinessCode `json:"code"`
	Message    string       `json:"message"`
	RawError   error        `json:"-"`
}

func (e *AppError) Error() string {
	if e.RawError != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.RawError)
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

func NewSMTPMailError(msg string, err error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       CodeSMTPError,
		Message:    msg,
		RawError:   err,
	}
}

func MapDBError(err error) *AppError {
	if errors.Is(err, sql.ErrNoRows) {
		return &AppError{
			HTTPStatus: http.StatusNotFound,
			Code:       CodeDBError,
			Message:    "The requested resource could not be found",
			RawError:   err,
		}
	}

	errStr := err.Error()
	if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "database is closed") {
		return &AppError{
			HTTPStatus: http.StatusInternalServerError,
			Code:       CodeDBError,
			Message:    "Our database system is currently unreachable. Please try again later.",
			RawError:   err,
		}
	}

	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       CodeDBError,
		Message:    "An unexpected error occurred while reading from our systems.",
		RawError:   err,
	}
}
