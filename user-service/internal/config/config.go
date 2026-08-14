package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
	DB  DatabaseConfig
}
type AppConfig struct {
	Name   string
	Port   string
	Env    string
	APIKey string
}
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	return &Config{
		App: AppConfig{
			Name:   os.Getenv("USER_APP_NAME"),
			Port:   os.Getenv("USER_APP_PORT"),
			Env:    os.Getenv("APP_ENV"),
			APIKey: os.Getenv("API_KEY"),
		},
		DB: DatabaseConfig{
			Host:     os.Getenv("USER_DB_HOST"),
			Port:     os.Getenv("USER_DB_PORT"),
			Username: os.Getenv("USER_DB_USERNAME"),
			Password: os.Getenv("USER_DB_PASSWORD"),
			Name:     os.Getenv("USER_DB_NAME"),
			SSLMode:  os.Getenv("USER_DB_SSLMODE"),
		},
	}
}
