package model

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;unique_index"`
	Email        string    `gorm:"type:varchar(255);not null;unique_index"`
	PasswordHash string    `gorm:"type:varchar(255);not null;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
