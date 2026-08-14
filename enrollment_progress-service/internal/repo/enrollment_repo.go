package repo

import (
	"context"
	"main/enrollment_progress-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnrollmentRepo interface {
	Create(ctx context.Context, enrollment *model.Enrollment) error
	GetByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (*model.Enrollment, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Enrollment, error)
	Update(ctx context.Context, enrollment *model.Enrollment) error
}
type enrollmentRepo struct {
	db *gorm.DB
}

func NewEnrollmentRepo(db *gorm.DB) EnrollmentRepo {
	return &enrollmentRepo{
		db: db,
	}
}
func (r *enrollmentRepo) Create(ctx context.Context, enrollment *model.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}
func (r *enrollmentRepo) GetByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (*model.Enrollment, error) {
	var enrollment model.Enrollment
	err := r.db.WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}
func (r *enrollmentRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Enrollment, error) {
	var enrollments []model.Enrollment
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&enrollments).Error
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}
func (r *enrollmentRepo) Update(ctx context.Context, enrollment *model.Enrollment) error {
	return r.db.WithContext(ctx).Save(enrollment).Error
}
