package model

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;primaryKey"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null;unique_index:idx_course_position"`
	Title     string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text;not null"`
	Position  int       `gorm:"type:int;not null;unique_index:idx_course_position"`
	IsPreview bool      `gorm:"type:boolean;not null;default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
