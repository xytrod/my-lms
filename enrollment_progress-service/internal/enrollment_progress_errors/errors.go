package enrollment_progress_errors

import (
	"errors"
)

var (
	ErrUnauthorized             = errors.New("Unauthorized")
	ErrAlreadyEnrolled          = errors.New("Already enrolled")
	ErrEnrollmentNotFound       = errors.New("Enrollment not found")
	ErrCourseNotFound           = errors.New("Course not found")
	ErrCourseNotActive          = errors.New("Course not active")
	ErrLessonNotFound           = errors.New("Lesson not found")
	ErrCourseServiceUnavailable = errors.New("Course service unavailable")
)
