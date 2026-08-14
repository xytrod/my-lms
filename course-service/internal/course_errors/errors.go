package course_errors

import "errors"

var (
	ErrCourseHasNoLessons  = errors.New("course has no lessons")
	ErrCourseNotFound      = errors.New("course not found")
	ErrForbidden           = errors.New("forbidden")
	ErrCourseAlreadyExists = errors.New("course already exists")
	ErrCourseArchived      = errors.New("course archived")
	ErrCourseLessonsLocked = errors.New("course lessons locked")
	ErrInvalidCourseTitle  = errors.New("invalid course title")
	ErrInvalidCourseID     = errors.New("invalid course id")
	ErrInvalidCourseLevel  = errors.New("invalid course level")

	ErrLessonNotFound        = errors.New("lesson not found")
	ErrInvalidLessonTitle    = errors.New("invalid lesson_title")
	ErrInvalidLessonID       = errors.New("invalid lesson_id")
	ErrInvalidLessonPosition = errors.New("invalid lesson_position")
	ErrInvalidLessonContent  = errors.New("invalid lesson_content")
)
