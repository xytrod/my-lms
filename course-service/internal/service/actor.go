package service

import "github.com/google/uuid"

type Actor struct {
	UserID uuid.UUID
	Role   string
}
