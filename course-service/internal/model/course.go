package model

import (
	"time"

	"github.com/google/uuid"
)

type CourseState string
type CourseDifficulty string

const (
	CourseStateDraft  CourseState = "Draft"
	CourseStateActive CourseState = "published"
	CourseStateClosed CourseState = "closed"
)
const (
	CourseDifficultyBeginner     CourseDifficulty = "beginner"
	CourseDifficultyIntermediate CourseDifficulty = "intermediate"
	CourseDifficultyAdvanced     CourseDifficulty = "advanced"
)

type Course struct {
	ID          uuid.UUID        `gorm:"primary_key;type:uuid;primaryKey"`
	TeacherID   uuid.UUID        `gorm:"type:uuid;not null;index"`
	Title       string           `gorm:"type:varchar(255);not null"`
	Description string           `gorm:"type:varchar(255);not null"`
	Level       CourseDifficulty `gorm:"type:varchar(255);not null;index"`
	State       CourseState      `gorm:"type:varchar(255);not null;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Lessons     []Lesson `gorm:"foreignkey:CourseID"`
}
