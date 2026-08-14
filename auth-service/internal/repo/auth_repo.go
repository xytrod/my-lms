package repo

import (
	"context"
	"main/auth-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(ctx context.Context, credential *model.Credential) error
	GetByEmail(ctx context.Context, email string) (*model.Credential, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) (*model.Credential, error)
	DeleteByUserId(ctx context.Context, userID uuid.UUID) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
type authRepo struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepo{
		db: db,
	}
}
func (r *authRepo) Create(ctx context.Context, credential *model.Credential) error {
	return r.db.Debug().WithContext(ctx).Create(credential).Error
}
func (r *authRepo) GetByEmail(ctx context.Context, email string) (*model.Credential, error) {
	var credential model.Credential
	err := r.db.Where("email = ?", email).First(&credential).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}
func (r *authRepo) GetByUserId(ctx context.Context, userID uuid.UUID) (*model.Credential, error) {
	var credential model.Credential
	err := r.db.Where("user_id = ?", userID).First(&credential).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}
func (r *authRepo) DeleteByUserId(ctx context.Context, userID uuid.UUID) error {
	res := r.db.Delete(&model.Credential{}, "user_id = ?", userID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *authRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Credential{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
