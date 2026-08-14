package service

import "errors"

var (
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrUserInActive            = errors.New("user is inactive")
	ErrCredentialAlreadyExists = errors.New("credential already exists")
	ErrRegistrationFailed      = errors.New("registration failed")
	ErrUserServiceUnavailable  = errors.New("user service unavailable")
	ErrInvalidRefreshToken     = errors.New("invalid refresh token")
)
