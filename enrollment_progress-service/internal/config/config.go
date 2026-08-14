package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App       AppConfig
	DB        DatabaseConfig
	JWT       JWTConfig
	CourseCfg CourseServiceConfig
	RabbitMQ  RabbitMQConfig
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
type CourseServiceConfig struct {
	URL string
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
			Name: os.Getenv("ENROLL_PROG_APP_NAME"),
			Port: os.Getenv("ENROLL_PROG_APP_PORT"),
			Env:  os.Getenv("APP_ENV"),
		},
		DB: DatabaseConfig{
			Host:     os.Getenv("ENROLL_PROG_DB_HOST"),
			Port:     os.Getenv("ENROLL_PROG_DB_PORT"),
			Username: os.Getenv("ENROLL_PROG_DB_USERNAME"),
			Password: os.Getenv("ENROLL_PROG_DB_PASSWORD"),
			Name:     os.Getenv("ENROLL_PROG_DB_NAME"),
			SSLMode:  os.Getenv("ENROLL_PROG_DB_SSLMODE"),
		},
		JWT: JWTConfig{
			AccessSecret: os.Getenv("ENROLL_PROG_ACCESS"),
			Issuer:       os.Getenv("ENROLL_PROG_ISSUER"),
		},
		CourseCfg: CourseServiceConfig{
			URL: os.Getenv("COURSE_ENROLL_SERVICE_URL"),
		},
		RabbitMQ: RabbitMQConfig{
			URL: os.Getenv("RABBITMQ_URL"),
		},
	}
	return cfg, nil
}
