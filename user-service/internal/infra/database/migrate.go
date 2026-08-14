package database

import (
	"main/user-service/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
	)
}
