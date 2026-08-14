package database

import (
	"main/auth-service/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Credential{}, &model.RefreshSession{})
}
