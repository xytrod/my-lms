package model

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	TokenID     string     `gorm:"type:varchar(70);not null;unique_index"`
	HashedToken string     `gorm:"type:varchar(70);not null;unique_index"`
	ExpiresAt   time.Time  `gorm:"not null;index"`
	RevokedAt   *time.Time `gorm:"index"`
	CreatedAt   time.Time
}
