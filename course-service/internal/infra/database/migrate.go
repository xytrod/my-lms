package database

import (
	"main/course-service/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Course{}, &model.Lesson{})
}
