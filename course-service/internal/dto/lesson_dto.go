package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateLessonRequest struct {
	Title     string `json:"title" validate:"required,min=10,max=25"`
	Content   string `json:"content" validate:"required,min=1"`
	Position  int    `json:"position" validate:"required,gte=1"`
	IsPreview bool   `json:"is_preview"`
}
type UpdateLessonRequest struct {
	Title     *string `json:"title" validate:"omitempty,min=10,max=25"`
	Content   *string `json:"content" validate:"omitempty,min=1"`
	Position  *int    `json:"position" validate:"omitempty,gte=1"`
	IsPreview *bool   `json:"is_preview"`
}
type LessonResponse struct {
	ID        uuid.UUID `json:"id"`
	CourseID  uuid.UUID `json:"course_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Position  int       `json:"position"`
	IsPreview bool      `json:"is_preview"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
