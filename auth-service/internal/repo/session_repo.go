package repo

import (
	"context"
	"main/auth-service/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepo interface {
	Create(ctx context.Context, session *model.RefreshSession) error
	GetByToken(ctx context.Context, token string) (*model.RefreshSession, error)
	Revoke(ctx context.Context, token string) error
	RevokeByID(ctx context.Context, userID uuid.UUID) error
}
type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepo {
	return &sessionRepository{
		db: db,
	}
}
func (r *sessionRepository) Create(ctx context.Context, session *model.RefreshSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}
func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*model.RefreshSession, error) {
	var session model.RefreshSession
	err := r.db.WithContext(ctx).First(&session, "token_id = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}
func (r *sessionRepository) Revoke(ctx context.Context, token string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&model.RefreshSession{}).Where("token_id = ? AND revoked_at IS NULL", token).Update("revoked_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *sessionRepository) RevokeByID(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&model.RefreshSession{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", now)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
