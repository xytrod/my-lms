package token

import "errors"

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrInvalidIssuer    = errors.New("invalid token issuer")
)
