package apperror

import (
	"net/http"
)

func New(status int, code string, message string, err error) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}
func BadRequest(code string, message string, err error) *AppError {
	return New(http.StatusBadRequest, code, message, err)
}
func Conflict(code string, message string, err error) *AppError {
	return New(http.StatusConflict, code, message, err)
}
func Internal(code string, message string, err error) *AppError {
	return New(http.StatusInternalServerError, code, message, err)
}
func ServiceUnavailable(code string, message string, err error) *AppError {
	return New(http.StatusServiceUnavailable, code, message, err)
}
func Unauthorized(code string, message string, err error) *AppError {
	return New(http.StatusUnauthorized, code, message, err)
}
