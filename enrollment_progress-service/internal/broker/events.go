package broker

import (
	"time"

	"github.com/google/uuid"
)

type UserEnrolledEvent struct {
	EnrollmentID uuid.UUID `json:"enrollment_id"`
	UserID       uuid.UUID `json:"user_id"`
	CourseID     uuid.UUID `json:"course_id"`
	OccurredAt   time.Time `json:"occurred_at"`
}
