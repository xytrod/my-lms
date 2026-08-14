package model

import (
	"time"

	"github.com/google/uuid"
)

type EnrollmentStatus string

const (
	EnrollmentStatusActive    EnrollmentStatus = "active"
	EnrollmentStatusCompleted EnrollmentStatus = "completed"
)

type Enrollment struct {
	ID          uuid.UUID        `gorm:"uuid;primaryKey"`
	UserID      uuid.UUID        `gorm:"type:uuid;not null;unique_index:idx_user_course"`
	CourseID    uuid.UUID        `gorm:"type:uuid;not null;unique_index:idx_user_course"`
	Status      EnrollmentStatus `gorm:"type:varchar(15);not null;index"`
	StartedAt   time.Time        `gorm:"not null"`
	CompletedAt *time.Time       `gorm:"default:null"`
	CreatedAt   time.Time        `gorm:"not null"`
	UpdatedAt   time.Time        `gorm:"not null"`
}
