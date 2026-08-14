package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App         AppConfig
	DB          DatabaseConfig
	JWT         JWTConfig
	UserService UserServiceCfg
}
type AppConfig struct {
	Name string
	Port string
	Env  string
}
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}
type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type UserServiceCfg struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning loading .env file")
	}
	accessTTL, err := time.ParseDuration(os.Getenv("AUTH_JWT_ACCESS_TTL"))
	if err != nil {
		return nil, err
	}
	refreshTTL, err := time.ParseDuration(os.Getenv("AUTH_JWT_REFRESH_TTL"))
	if err != nil {
		return nil, err
	}
	accessSecret := os.Getenv("AUTH_JWT_ACCESS")
	if accessSecret == "" {
		return nil, fmt.Errorf("access secret required")
	}
	refreshSecret := os.Getenv("AUTH_JWT_REFRESH")
	if refreshSecret == "" {
		return nil, fmt.Errorf("refresh secret required")
	}
	userServiceTimeout, err := time.ParseDuration(os.Getenv("AUTH_USER_TIMEOUT"))
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		App: AppConfig{
			Name: os.Getenv("AUTH_APP_NAME"),
			Port: os.Getenv("AUTH_APP_PORT"),
			Env:  os.Getenv("APP_ENV"),
		},
		DB: DatabaseConfig{
			Host:     os.Getenv("AUTH_DB_HOST"),
			Port:     os.Getenv("AUTH_DB_PORT"),
			Username: os.Getenv("AUTH_DB_USERNAME"),
			Password: os.Getenv("AUTH_DB_PASSWORD"),
			Name:     os.Getenv("AUTH_DB_NAME"),
			SSLMode:  os.Getenv("AUTH_DB_SSLMODE"),
		},
		JWT: JWTConfig{
			AccessSecret:  accessSecret,
			RefreshSecret: refreshSecret,
			AccessTTL:     accessTTL,
			RefreshTTL:    refreshTTL,
		},
		UserService: UserServiceCfg{
			BaseURL: os.Getenv("AUTH_USER_URL"),
			APIKey:  os.Getenv("AUTH_USER_API_KEY"),
			Timeout: userServiceTimeout,
		},
	}
	return cfg, nil
}
