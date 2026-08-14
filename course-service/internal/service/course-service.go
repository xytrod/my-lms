package service

import (
	"context"
	"errors"
	"main/course-service/internal/apperror"
	"main/course-service/internal/course_errors"
	"main/course-service/internal/dto"
	"main/course-service/internal/model"
	"main/course-service/internal/repo"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseService interface {
	ListLessons(ctx context.Context, courseID uuid.UUID) ([]model.Lesson, error)
	GetManagedCourse(ctx context.Context, actor Actor, id uuid.UUID) (*model.Course, error)
	ListManagedLessons(ctx context.Context, actor Actor, courseID uuid.UUID) ([]model.Lesson, error)
	Create(ctx context.Context, actor Actor, req dto.CourseRequest) (*model.Course, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Course, error)
	Update(ctx context.Context, actor Actor, id uuid.UUID, req dto.UpdateCourseRequest) (*model.Course, error)
	Publish(ctx context.Context, actor Actor, id uuid.UUID) (*model.Course, error)
	Archive(ctx context.Context, actor Actor, id uuid.UUID) (*model.Course, error)
	ListPublic(ctx context.Context, search string, limit, offset int) ([]model.Course, error)
	ListMyCourses(ctx context.Context, actor Actor, limit, offset int) ([]model.Course, error)
	CreateLesson(ctx context.Context, actor Actor, courseID uuid.UUID, req dto.CreateLessonRequest) (*model.Lesson, error)
	UpdateLesson(ctx context.Context, actor Actor, lessonID uuid.UUID, req dto.UpdateLessonRequest) (*model.Lesson, error)
	DeleteLesson(ctx context.Context, actor Actor, lessonID uuid.UUID) error
}
type courseService struct {
	CourseRepo repo.CourseRepository
	LessonRepo repo.LessonRepository
}

func NewCourseService(courseRepo repo.CourseRepository, lessonRepo repo.LessonRepository) CourseService {
	return &courseService{
		CourseRepo: courseRepo,
		LessonRepo: lessonRepo,
	}
}

func (s *courseService) GetManagedCourse(ctx context.Context, actor Actor, id uuid.UUID) (*model.Course, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("user_unauthorized", "authentification is required", course_errors.ErrForbidden)
	}
	if id == uuid.Nil {
		return nil, apperror.BadRequest("invalid_course_id", "invalid course id", course_errors.ErrInvalidCourseID)
	}
	course, err := s.CourseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
		return nil, apperror.Internal("course_get_failed", "course could not be loaded", err)
	}
	if !IsManagingAble(actor, course) {
		return nil, apperror.Forbidden("access_forbidden", "you are not allowed to manage this course", course_errors.ErrForbidden)
	}
	return course, nil
}

func (s *courseService) ListManagedLessons(ctx context.Context, actor Actor, courseID uuid.UUID) ([]model.Lesson, error) {
	if _, err := s.GetManagedCourse(ctx, actor, courseID); err != nil {
		return nil, err
	}
	lessons, err := s.LessonRepo.ListByCourseID(ctx, courseID)
	if err != nil {
		return nil, apperror.Internal("lesson_storage_error", "lesson service is unavailable now", err)
	}
	return lessons, nil
}

func lessonsLockedError() error {
	return apperror.Conflict(
		"course_lessons_locked",
		"lesson structure can only be changed while the course is a draft",
		course_errors.ErrCourseLessonsLocked,
	)
}
func (s *courseService) Create(ctx context.Context, actor Actor, req dto.CourseRequest) (*model.Course, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("user_unauthorized", "authentification required", course_errors.ErrForbidden)
	}
	title := strings.TrimSpace(req.Title)
	desc := strings.TrimSpace(req.Description)
	if title == "" {
		return nil, apperror.BadRequest("invalid_course_title", "course title is invalid", course_errors.ErrInvalidCourseTitle)
	}
	level := model.CourseDifficulty(req.Level)
	switch level {
	case model.CourseDifficultyBeginner,
		model.CourseDifficultyIntermediate,
		model.CourseDifficultyAdvanced:
	default:
		return nil, apperror.BadRequest("invalid_course_level", "course level is invalid", course_errors.ErrInvalidCourseLevel)
	}
	course := &model.Course{
		ID:          uuid.New(),
		TeacherID:   actor.UserID,
		Title:       title,
		Description: desc,
		Level:       level,
		State:       model.CourseStateDraft,
	}
	if err := s.CourseRepo.Create(ctx, course); err != nil {
		return nil, apperror.Internal("course_create_failed", "course isnt created", err)
	}
	return course, nil
}
func (s *courseService) GetByID(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	if id == uuid.Nil {
		return nil, apperror.BadRequest("invalid_course_id", "invalid course id", course_errors.ErrInvalidCourseID)
	}
	course, err := s.CourseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
		return nil, apperror.Internal("course_get_failed", "course isnt found", err)
	}
	if course.State != model.CourseStateActive {
		return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
	}
	return course, nil
}
func (s *courseService) Update(ctx context.Context, actor Actor, id uuid.UUID, req dto.UpdateCourseRequest) (*model.Course, error) {
	course, err := s.CourseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
	}
	if !IsManagingAble(actor, course) {
		return nil, apperror.Forbidden("access_forbidden", "you are not allowed to do this action", course_errors.ErrForbidden)
	}
	if course.State == model.CourseStateClosed {
		return nil, apperror.Conflict("course_archived", "course archived", course_errors.ErrCourseArchived)
	}
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, apperror.BadRequest("invalid_course_title", "invalid title", course_errors.ErrInvalidCourseTitle)
		}
		course.Title = title
	}
	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		if desc == "" {
			return nil, errors.New("description is required")
		}
		course.Description = desc
	}
	if req.Level != nil {
		level := model.CourseDifficulty(strings.TrimSpace(*req.Level))
		switch level {
		case model.CourseDifficultyBeginner,
			model.CourseDifficultyIntermediate,
			model.CourseDifficultyAdvanced:
		default:
			return nil, apperror.BadRequest("invalid_course_level", "invalid level", course_errors.ErrInvalidCourseLevel)
		}
		course.Level = level
	}
	if err := s.CourseRepo.Update(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}
func (s *courseService) Publish(ctx context.Context, actor Actor, id uuid.UUID) (*model.Course, error) {
	course, err := s.CourseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
		return nil, err
	}
	if !IsManagingAble(actor, course) {
		return nil, apperror.Forbidden("access_forbidden", "you are not allowed to do this action", course_errors.ErrForbidden)
	}
	if course.State == model.CourseStateClosed {
		return nil, course_errors.ErrCourseArchived
	}
	if course.State == model.CourseStateActive {
		return nil, course_errors.ErrCourseAlreadyExists
	}
	lessonCount, err := s.LessonRepo.CountByCourseID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("lesson_count_failed", "lesson_count", err)
	}
	if lessonCount == 0 {
		return nil, apperror.Internal("course_hasnt_lessons", "course_hasnt_lessons", course_errors.ErrCourseHasNoLessons)
	}
	course.State = model.CourseStateActive
	if err := s.CourseRepo.Update(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}
func (s *courseService) Archive(ctx context.Context, actor Actor, id uuid.UUID) (*model.Course, error) {
	course, err := s.CourseRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
		return nil, err
	}
	if !IsManagingAble(actor, course) {
		return nil, apperror.Forbidden("access_forbidden", "you are not allowed to do this action", course_errors.ErrForbidden)
	}
	if course.State == model.CourseStateClosed {
		return course, nil
	}
	course.State = model.CourseStateClosed
	if err := s.CourseRepo.Update(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}
func (s *courseService) CreateLesson(ctx context.Context, actor Actor, courseID uuid.UUID, req dto.CreateLessonRequest) (*model.Lesson, error) {
	course, err := s.CourseRepo.GetByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
		return nil, apperror.Internal("course_storage_error", "course service is unavailable now", err)
	}
	if !IsManagingAble(actor, course) {
		return nil, apperror.Forbidden("course_forbidden", "you are not allowed to do this action", course_errors.ErrForbidden)
	}
	if course.State != model.CourseStateDraft {
		return nil, lessonsLockedError()
	}
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	lesson := &model.Lesson{
		ID:        uuid.New(),
		CourseID:  courseID,
		Title:     title,
		Content:   content,
		Position:  req.Position,
		IsPreview: req.IsPreview,
	}
	if err := s.LessonRepo.Create(ctx, lesson); err != nil {
		return nil, apperror.Internal("lesson_storage_error", "lesson_storage service is unavailable now", err)
	}
	return lesson, nil
}
func (s *courseService) UpdateLesson(ctx context.Context, actor Actor, lessonID uuid.UUID, req dto.UpdateLessonRequest) (*model.Lesson, error) {
	if lessonID == uuid.Nil {
		return nil, apperror.BadRequest("lesson_id", "lesson_id is required", course_errors.ErrInvalidLessonID)
	}
	lesson, err := s.LessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("lesson_not_found", "lesson not found", course_errors.ErrLessonNotFound)
		}
		return nil, apperror.Internal("lesson_storage_error", "lesson service is unavailable now", err)
	}
	course, err := s.CourseRepo.GetByID(ctx, lesson.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
		return nil, apperror.Internal("course_storage_error", "course service is unavailable now", err)
	}
	if !IsManagingAble(actor, course) {
		return nil, apperror.Forbidden("course_forbidden", "you are not allowed to do this action", course_errors.ErrForbidden)
	}
	if course.State != model.CourseStateDraft {
		return nil, lessonsLockedError()
	}
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, apperror.BadRequest("title", "title is required", course_errors.ErrInvalidLessonTitle)
		}
		lesson.Title = title
	}
	if req.Content != nil {
		content := strings.TrimSpace(*req.Content)
		if content == "" {
			return nil, apperror.BadRequest("content", "content is required", course_errors.ErrInvalidLessonContent)
		}
		lesson.Content = content
	}
	if req.Position != nil {
		if *req.Position < 0 {
			return nil, apperror.BadRequest("position", "position is required", course_errors.ErrInvalidLessonPosition)
		}
		lesson.Position = *req.Position
	}
	if req.IsPreview != nil {
		lesson.IsPreview = *req.IsPreview
	}
	if err := s.LessonRepo.Update(ctx, lesson); err != nil {
		return nil, apperror.Internal("lesson_storage_error", "lesson service is unavailable now", err)
	}
	return lesson, nil
}
func (s *courseService) DeleteLesson(ctx context.Context, actor Actor, lessonID uuid.UUID) error {
	if lessonID == uuid.Nil {
		return apperror.BadRequest("lesson_id", "lesson_id is required", course_errors.ErrInvalidLessonID)
	}
	lesson, err := s.LessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("lesson_not_found", "lesson not found", course_errors.ErrLessonNotFound)
		}
		return apperror.Internal("lesson_storage_error", "lesson service is unavailable now", err)
	}
	course, err := s.CourseRepo.GetByID(ctx, lesson.CourseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
		return apperror.Internal("course_storage_error", "course service is unavailable now", err)
	}
	if !IsManagingAble(actor, course) {
		return apperror.Forbidden("course_forbidden", "you are not allowed to do this action", course_errors.ErrForbidden)
	}
	if course.State != model.CourseStateDraft {
		return lessonsLockedError()
	}
	if err := s.LessonRepo.Delete(ctx, lesson.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("lesson_not_found", "lesson not found", course_errors.ErrLessonNotFound)
		}
		return apperror.Internal("lesson_storage_error", "lesson service is unavailable now", err)
	}
	return nil
}
func (s *courseService) ListLessons(ctx context.Context, courseID uuid.UUID) ([]model.Lesson, error) {
	if courseID == uuid.Nil {
		return nil, apperror.BadRequest("course_id", "course_id is required", course_errors.ErrInvalidCourseID)
	}
	course, err := s.CourseRepo.GetByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
		}
		return nil, apperror.Internal("course_storage_error", "course service is unavailable now", err)
	}
	if course.State != model.CourseStateActive {
		return nil, apperror.NotFound("course_not_found", "course not found", course_errors.ErrCourseNotFound)
	}
	lessons, err := s.LessonRepo.ListByCourseID(ctx, courseID)
	if err != nil {
		return nil, apperror.Internal("lesson_storage_error", "lesson service is unavailable now", err)
	}
	return lessons, nil
}
func (s *courseService) ListPublic(ctx context.Context, search string, limit, offset int) ([]model.Course, error) {
	search = strings.TrimSpace(search)
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	courses, err := s.CourseRepo.ListPublic(ctx, search, limit, offset)
	if err != nil {
		return nil, apperror.Internal("course_storage_error", "course loading is ended", err)
	}
	return courses, nil
}
func (s *courseService) ListMyCourses(ctx context.Context, actor Actor, limit, offset int) ([]model.Course, error) {
	if actor.UserID == uuid.Nil {
		return nil, apperror.Unauthorized("user_unauthorized", "authentification is required", course_errors.ErrForbidden)
	}
	courses, err := s.CourseRepo.ListByTeacher(ctx, actor.UserID, limit, offset)
	if err != nil {
		return nil, apperror.Internal("course_storage_error", "course loading is ended", err)
	}
	return courses, nil
}
func IsManagingAble(actor Actor, course *model.Course) bool {
	if actor.Role == "admin" {
		return true
	}
	return actor.UserID == course.TeacherID
}

func ToCourseDTO(course *model.Course) dto.CourseResponse {
	return dto.CourseResponse{
		ID:          course.ID,
		TeacherID:   course.TeacherID,
		Title:       course.Title,
		Description: course.Description,
		Level:       string(course.Level),
		State:       string(course.State),
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}
}
func ToLessonDTO(lesson *model.Lesson) dto.LessonResponse {
	return dto.LessonResponse{
		ID:        lesson.ID,
		CourseID:  lesson.CourseID,
		Title:     lesson.Title,
		Content:   lesson.Content,
		Position:  lesson.Position,
		IsPreview: lesson.IsPreview,
		CreatedAt: lesson.CreatedAt,
		UpdatedAt: lesson.UpdatedAt,
	}
}
func ToCourseDTOList(courses []model.Course) []dto.CourseResponse {
	res := make([]dto.CourseResponse, 0, len(courses))
	for i := range courses {
		res = append(res, ToCourseDTO(&courses[i]))
	}
	return res
}
