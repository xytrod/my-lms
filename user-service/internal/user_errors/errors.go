package user_errors

import (
	"errors"
)

var (
	ErrEmailRequired         = errors.New("email is required")
	ErrUsernameRequired      = errors.New("username is required")
	ErrUsernameAlreadyExists = errors.New("your username already exists")
	ErrPasswordRequired      = errors.New("password is required")
	ErrFirstNameRequired     = errors.New("your first name is required")
	ErrLastNameRequired      = errors.New("your last name is required")
	ErrInvalidEmail          = errors.New("your email is invalid")
	ErrEmailAlreadyExists    = errors.New("your email already exists")
	ErrInvalidUserID         = errors.New("invalid user ID")
	ErrUserNotFound          = errors.New("user not found")
)
