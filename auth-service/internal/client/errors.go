package client

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserService       = errors.New("user service course_errors")
	ErrInvalidResponse   = errors.New("invalid user service response")
)
