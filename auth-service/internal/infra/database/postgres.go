package database

import (
	"fmt"
	"main/auth-service/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port =%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("course_errors connecting to postgres: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("course_errors connecting to sql: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("pinging sql: %w", err)
	}
	return db, nil
}
