package token

import "errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
	ErrInvalidTokenType = errors.New("invalid token type")
)
