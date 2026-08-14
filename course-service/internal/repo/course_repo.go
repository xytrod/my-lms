package repo

import (
	"context"
	"main/course-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(ctx context.Context, course *model.Course) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Course, error)
	Update(ctx context.Context, course *model.Course) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListPublished(ctx context.Context, limit, offset int) ([]*model.Course, error)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]model.Course, error)
	ListPublic(ctx context.Context, search string, limit, offset int) ([]model.Course, error)
}
type courseRepo struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepo{
		db: db,
	}
}
func (r *courseRepo) Create(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}
func (r *courseRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Course, error) {
	var course model.Course
	err := r.db.WithContext(ctx).First(&course, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}
func (r *courseRepo) Update(ctx context.Context, course *model.Course) error {
	res := r.db.Debug().WithContext(ctx).Save(course)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *courseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Course{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *courseRepo) ListPublished(ctx context.Context, limit, offset int) ([]*model.Course, error) {
	var courses []*model.Course
	err := r.db.WithContext(ctx).Where("status = ?", model.CourseStateActive).Order("created_at DESC").Limit(limit).Offset(offset).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
func (r *courseRepo) ListByTeacher(ctx context.Context, teacherID uuid.UUID, limit, offset int) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.WithContext(ctx).Where("teacher_id = ?", teacherID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
func (r *courseRepo) ListPublic(ctx context.Context, search string, limit, offset int) ([]model.Course, error) {
	var courses []model.Course
	query := r.db.WithContext(ctx).Where("state = ?", model.CourseStateActive)
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
