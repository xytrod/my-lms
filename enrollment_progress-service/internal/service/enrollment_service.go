package service

import (
	"context"
	"errors"
	"log"
	"main/enrollment_progress-service/internal/apperror"
	"main/enrollment_progress-service/internal/broker"
	"main/enrollment_progress-service/internal/client"
	"main/enrollment_progress-service/internal/dto"
	enroll_prog__errors "main/enrollment_progress-service/internal/enrollment_progress_errors"
	"main/enrollment_progress-service/internal/model"
	"main/enrollment_progress-service/internal/repo"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentService interface {
	Enroll(ctx context.Context, userID, courseID uuid.UUID) (*model.Enrollment, error)
	CompletedLesson(ctx context.Context, userID, courseID, lessonID uuid.UUID) error
	GetProgress(ctx context.Context, userID, courseID uuid.UUID) (*dto.ProgressResponse, error)
	ListMy(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Enrollment, error)
}
type enrollmentService struct {
	enrollmentRepo repo.EnrollmentRepo
	progressRepo   repo.ProgressRepo
	courseClient   client.CourseClient
	publisher      *broker.Publisher
}

func NewEnrollmentService(enrollmentRepo repo.EnrollmentRepo, progressRepo repo.ProgressRepo, courseClient client.CourseClient, publisher *broker.Publisher) *enrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		progressRepo:   progressRepo,
		courseClient:   courseClient,
		publisher:      publisher,
	}
}
func (s *enrollmentService) Enroll(ctx context.Context, userID, courseID uuid.UUID) (*model.Enrollment, error) {
	if userID == uuid.Nil {
		return nil, apperror.Unauthorized("User_ID_required", "authentification required", enroll_prog__errors.ErrUnauthorized)
	}
	log.Printf(
		"ENROLL: userID=%s courseID=%s",
		userID,
		courseID,
	)
	course, err := s.courseClient.GetCourseByCourseID(ctx, courseID)
	if err != nil {
		log.Printf("ENROLL: course client error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Internal("course_not_found", "course not found", enroll_prog__errors.ErrCourseNotFound)
		}
		if errors.Is(err, enroll_prog__errors.ErrCourseServiceUnavailable) {
			return nil, apperror.ServiceUnavailable("course_service_unavailable", "course service is unavailable now", err)
		}
		return nil, apperror.Internal("enrollment_service_error", "enrollment service error", err)
	}

	if course.State != "published" {
		return nil, apperror.Conflict("course_not_active", "course is not available now", enroll_prog__errors.ErrCourseNotActive)
	}
	exists, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err == nil && exists != nil {
		return nil, apperror.Conflict("already_enrolled", "enrollment could not be processed", enroll_prog__errors.ErrAlreadyEnrolled)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.Internal("enrollment_error", "enrollment could not be processed", err)
	}
	enrollment := &model.Enrollment{
		ID:        uuid.New(),
		UserID:    userID,
		CourseID:  courseID,
		Status:    model.EnrollmentStatusActive,
		StartedAt: time.Now().UTC(),
	}
	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, apperror.Internal("enrollment_error", "enrollment could not be processed", err)
	}
	event := broker.UserEnrolledEvent{
		EnrollmentID: enrollment.ID,
		UserID:       enrollment.UserID,
		CourseID:     enrollment.CourseID,
		OccurredAt:   time.Now().UTC(),
	}
	log.Printf("publishing user.enrolled user=%s course=%s", event.UserID, event.CourseID)
	if err := s.publisher.Publish(ctx, "user.enrolled", event); err != nil {
		log.Printf("failed to publish user enrolled event: %v", err)
	}
	return enrollment, nil
}
func (s *enrollmentService) CompletedLesson(ctx context.Context, userID, courseID, lessonID uuid.UUID) error {
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("enrollment_not_found", "user is not enrolled in the course", enroll_prog__errors.ErrEnrollmentNotFound)
		}
		return apperror.Internal("enrollment_error", "enrollment could not be processed", err)
	}
	lessons, err := s.courseClient.GetLessonsByCourseID(ctx, courseID)
	if err != nil {
		return err
	}
	found_lesson := false
	for _, lesson := range lessons {
		if lesson.ID == lessonID {
			found_lesson = true
			break
		}
	}
	if !found_lesson {
		return apperror.NotFound("lesson_not_found", "lesson not found in this course", enroll_prog__errors.ErrLessonNotFound)
	}
	exists, err := s.progressRepo.Exists(ctx, enrollment.ID, lessonID)
	if err != nil {
		return apperror.Internal("enrollment_error", "enrollment could not be processed", err)
	}
	if exists {
		return nil
	}
	progress := &model.LessonProgress{
		ID:           uuid.New(),
		EnrollmentID: enrollment.ID,
		LessonID:     lessonID,
		CompletedAt:  time.Now().UTC(),
	}
	if err := s.progressRepo.Create(ctx, progress); err != nil {
		return apperror.Internal("enrollment_error", "enrollment could not be processed", err)
	}
	completed, err := s.progressRepo.CountCompleted(ctx, enrollment.ID)
	if err != nil {
		return apperror.Internal("progress_error", "progress could not be processed", err)
	}
	if len(lessons) > 0 && completed == int64(len(lessons)) {
		now := time.Now().UTC()
		enrollment.Status = model.EnrollmentStatusCompleted
		enrollment.CompletedAt = &now
		if err := s.enrollmentRepo.Update(ctx, enrollment); err != nil {
			return apperror.Internal("enrollment_error", "enrollment could not be processed", err)
		}
	}
	return nil
}
func (s *enrollmentService) GetProgress(ctx context.Context, userID, courseID uuid.UUID) (*dto.ProgressResponse, error) {
	var percentage float64
	enrollment, err := s.enrollmentRepo.GetByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	lessons, err := s.courseClient.GetLessonsByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	completed, err := s.progressRepo.CountCompleted(ctx, enrollment.ID)
	if err != nil {
		return nil, err
	}
	total := len(lessons)
	if total > 0 {
		percentage = float64(completed) / float64(total) * 100
	}
	return &dto.ProgressResponse{
		Completed:  completed,
		Total:      int64(total),
		Percentage: percentage,
	}, nil
}
func (s *enrollmentService) ListMy(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Enrollment, error) {
	if userID == uuid.Nil {
		return nil, apperror.Unauthorized("user_id", "user_id is required", enroll_prog__errors.ErrUnauthorized)
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	enrollments, err := s.enrollmentRepo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, apperror.Internal("enrollment_error", "enrollment could not be processed", err)
	}
	return enrollments, nil
}
