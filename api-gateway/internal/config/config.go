package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	AuthURL       string
	CourseURL     string
	EnrollmentURL string
	UserURL       string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Port:          os.Getenv("API_GATEWAY_PORT"),
		AuthURL:       os.Getenv("AUTH_SERVICE_URL"),
		CourseURL:     os.Getenv("COURSE_SERVICE_URL"),
		EnrollmentURL: os.Getenv("ENROLLMENT_PROG_SERVICE_URL"),
		UserURL:       os.Getenv("USER_SERVICE_URL"),
	}
	log.Printf(
		"CONFIG port=%q auth=%q user=%q course=%q enrollment=%q",
		cfg.Port,
		cfg.AuthURL,
		cfg.UserURL,
		cfg.CourseURL,
		cfg.EnrollmentURL,
	)
	if cfg.Port == "" {
		return nil, fmt.Errorf("API_GATEWAY_PORT is required")
	}
	if cfg.AuthURL == "" || cfg.CourseURL == "" || cfg.EnrollmentURL == "" || cfg.UserURL == "" {
		return nil, fmt.Errorf("invalid configuration")
	}
	return cfg, nil
}
