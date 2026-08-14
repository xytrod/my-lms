package token

import (
	"errors"
	"fmt"
	"main/auth-service/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type RotationRefresh struct {
	token     string
	TokenID   string
	ExpiresAt time.Time
}

func NewJWTManager(cfg config.JWTConfig) Manager {
	return &JWTManager{
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessTTL:     cfg.AccessTTL,
		refreshTTL:    cfg.RefreshTTL,
	}
}
func (m *JWTManager) generateAccessToken(userID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		UserID: userID,
		Role:   role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    "auth-service",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.accessSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return signedToken, nil
}
func (m *JWTManager) generateRefreshToken(userID uuid.UUID) (*RotationRefresh, error) {
	now := time.Now()
	tokenID := uuid.NewString()
	expiresAt := now.Add(m.refreshTTL)
	claims := AccessClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    "auth-service",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        tokenID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return &RotationRefresh{
		token:     signedToken,
		TokenID:   tokenID,
		ExpiresAt: expiresAt,
	}, nil
}
func (m *JWTManager) GeneratePair(userID uuid.UUID, role string) (*Token, error) {
	accessToken, err := m.generateAccessToken(userID, role)
	if err != nil {
		return nil, err
	}
	refreshToken, err := m.generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	return &Token{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken.token,
		RefreshTokenID:   refreshToken.TokenID,
		RefreshExpiresAt: refreshToken.ExpiresAt,
	}, nil
}

func (m *JWTManager) ParseAccessToken(rawToken string) (*AccessClaims, error) {
	if rawToken == "" {
		return nil, ErrInvalidToken
	}
	claims := &AccessClaims{}
	parsed, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}

		return m.accessSecret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("auth-service"),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("failed to parse access token: %w", err)
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
func (m *JWTManager) ParseRefreshToken(rawToken string) (*RefreshClaims, error) {
	if rawToken == "" {
		return nil, ErrInvalidToken
	}
	claims := &RefreshClaims{}
	parsed, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return m.refreshSecret, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("auth-service"),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}
	if !parsed.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Type != "refresh" {
		return nil, ErrInvalidTokenType
	}
	if claims.UserID == uuid.Nil {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
