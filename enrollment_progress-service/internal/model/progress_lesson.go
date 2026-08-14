package model

import (
	"time"

	"github.com/google/uuid"
)

type LessonProgress struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	EnrollmentID uuid.UUID `gorm:"type:uuid;not null;unique_index:idx_enrollment_lesson"`
	LessonID     uuid.UUID `gorm:"type:uuid;not null;unique_index:idx_enrollment_lesson"`
	CompletedAt  time.Time `gorm:"not null"`
	CreatedAt    time.Time
}
