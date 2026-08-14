package dto

import (
	"time"

	"github.com/google/uuid"
)

type CourseRequest struct {
	Title       string `json:"title" validate:"required,min=20,max=60"`
	Description string `json:"description" validate:"required,min=20,max=200"`
	Level       string `json:"level" validate:"required,oneof=beginner intermediate advanced"`
}
type CourseResponse struct {
	ID          uuid.UUID `json:"id"`
	TeacherID   uuid.UUID `json:"teacher_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Level       string    `json:"level"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type UpdateCourseRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=20,max=60"`
	Description *string `json:"description" validate:"omitempty,min=20,max=200"`
	Level       *string `json:"level" validate:"omitempty,oneof=beginner intermediate advanced"`
}
