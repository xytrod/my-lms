package hashing

import "errors"

var (
	ErrPasswordMismatch = errors.New("password mismatch")
	ErrPasswordTooLong  = errors.New("password too long")
)
