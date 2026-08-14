package repository

import (
	"context"
	"main/user-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit int, offset int) ([]model.User, error)
}
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	res := r.db.Debug().WithContext(ctx).Save(user)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *UserRepository) List(ctx context.Context, limit int, offset int) ([]model.User, error) {
	var users []model.User
	err := r.db.WithContext(ctx).Order("created_at desc").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
