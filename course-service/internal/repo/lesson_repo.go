package repo

import (
	"context"
	"main/course-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *model.Lesson) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Lesson, error)
	ListByCourseID(ctx context.Context, courseID uuid.UUID) ([]model.Lesson, error)
	Update(ctx context.Context, lesson *model.Lesson) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByCourseID(ctx context.Context, courseID uuid.UUID) (int64, error)
}
type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{
		db: db,
	}
}
func (r *lessonRepository) Create(ctx context.Context, lesson *model.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}
func (r *lessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Lesson, error) {
	var lesson model.Lesson
	err := r.db.WithContext(ctx).First(&lesson, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}
func (r *lessonRepository) ListByCourseID(ctx context.Context, courseID uuid.UUID) ([]model.Lesson, error) {
	var lessons []model.Lesson
	err := r.db.WithContext(ctx).Find(&lessons, "course_id = ?", courseID).Order("position ASC").Error
	if err != nil {
		return nil, err
	}
	return lessons, nil
}
func (r *lessonRepository) Update(ctx context.Context, lesson *model.Lesson) error {
	res := r.db.WithContext(ctx).Save(lesson)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *lessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Lesson{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *lessonRepository) CountByCourseID(ctx context.Context, courseID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Lesson{}).Where("course_id = ?", courseID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
