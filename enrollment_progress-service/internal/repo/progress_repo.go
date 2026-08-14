package repo

import (
	"context"
	"main/enrollment_progress-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProgressRepo interface {
	Create(ctx context.Context, progress *model.LessonProgress) error
	Exists(ctx context.Context, enrollmentID, lessonID uuid.UUID) (bool, error)
	CountCompleted(ctx context.Context, enrollmentID uuid.UUID) (int64, error)
}
type progressRepo struct {
	db *gorm.DB
}

func NewProgressRepo(db *gorm.DB) ProgressRepo {
	return &progressRepo{
		db: db,
	}
}
func (r *progressRepo) Create(ctx context.Context, progress *model.LessonProgress) error {
	return r.db.WithContext(ctx).Create(progress).Error
}
func (r *progressRepo) Exists(ctx context.Context, enrollmentID, lessonID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LessonProgress{}).Where("enrollment_id = ? AND lesson_id = ?", enrollmentID, lessonID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (r *progressRepo) CountCompleted(ctx context.Context, enrollmentID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LessonProgress{}).Where("enrollment_id = ?", enrollmentID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
