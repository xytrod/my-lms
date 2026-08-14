package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	DB       DatabaseConfig
	JWT      JWTConfig
	RabbitMQ RabbitMQConfig
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
	AccessSecret string
	Issuer       string
}
type RabbitMQConfig struct {
	URL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning loading .env file")
	}
	cfg := &Config{
		App: AppConfig{
			Name: os.Getenv("COURSE_APP_NAME"),
			Port: os.Getenv("COURSE_APP_PORT"),
			Env:  os.Getenv("APP_ENV"),
		},
		DB: DatabaseConfig{
			Host:     os.Getenv("COURSE_DB_HOST"),
			Port:     os.Getenv("COURSE_DB_PORT"),
			Username: os.Getenv("COURSE_DB_USERNAME"),
			Password: os.Getenv("COURSE_DB_PASSWORD"),
			Name:     os.Getenv("COURSE_DB_NAME"),
			SSLMode:  os.Getenv("COURSE_DB_SSLMODE"),
		},
		JWT: JWTConfig{
			AccessSecret: os.Getenv("COURSE_JWT_ACCESS"),
			Issuer:       os.Getenv("COURSE_ISSUER"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: os.Getenv("RABBITMQ_URL"),
		},
	}
	return cfg, nil
}
