package token

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	AccessToken      string
	RefreshToken     string
	RefreshTokenID   string
	RefreshExpiresAt time.Time
}
type Manager interface {
	GeneratePair(userID uuid.UUID, role string) (*Token, error)
	ParseAccessToken(rawToken string) (*AccessClaims, error)
	ParseRefreshToken(rawToken string) (*RefreshClaims, error)
}
