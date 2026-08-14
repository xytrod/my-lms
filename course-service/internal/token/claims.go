package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	Type   string    `json:"type"`
	jwt.RegisteredClaims
}
