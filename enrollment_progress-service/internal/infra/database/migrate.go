package database

import (
	"main/enrollment_progress-service/internal/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Enrollment{}, &model.LessonProgress{})
}
