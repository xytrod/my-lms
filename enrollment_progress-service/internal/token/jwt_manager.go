package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	accessSecret []byte
	issuer       string
}

var _ Manager = (*JWTManager)(nil)

func NewJWTManager(accessSecret string, issuer string) Manager {
	return &JWTManager{
		accessSecret: []byte(accessSecret),
		issuer:       issuer,
	}
}
func (m *JWTManager) ParseAccessToken(token string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.accessSecret, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired())
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
