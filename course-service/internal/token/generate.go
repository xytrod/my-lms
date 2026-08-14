package token

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager interface {
	ParseAccessToken(token string) (*AccessClaims, error)
}

func (m *JWTManager) ParseAccessToken(token string) (*AccessClaims, error) {
	if strings.TrimSpace(token) == "" {
		return nil, ErrInvalidToken
	}
	claims := &AccessClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.accessSecret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, err
	}
	if !parsed.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Type != "access" {
		return nil, ErrInvalidTokenType
	}
	if claims.UserID == uuid.Nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
