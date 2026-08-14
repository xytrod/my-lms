package main

import (
	"log"
	"main/enrollment_progress-service/internal/broker"
	"main/enrollment_progress-service/internal/client"
	"main/enrollment_progress-service/internal/config"
	"main/enrollment_progress-service/internal/handler"
	"main/enrollment_progress-service/internal/infra/database"
	"main/enrollment_progress-service/internal/middleware"
	"main/enrollment_progress-service/internal/repo"
	"main/enrollment_progress-service/internal/routes"
	"main/enrollment_progress-service/internal/service"
	"main/enrollment_progress-service/internal/token"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load", err)
	}
	log.Printf("RMQ URL = %q", cfg.RabbitMQ.URL)
	rabbit, err := broker.NewRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("RMQ conn", err)
	}
	defer func() {
		if err := rabbit.Close(); err != nil {
			log.Printf("failed to close RabbitMQ connection: %v", err)
		}
	}()
	publisher := broker.NewPublisher(rabbit)
	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatal("conn", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal("database migration", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("db", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatal("close", err)
		}
	}()
	enrollmentRepo := repo.NewEnrollmentRepo(db)
	progressRepo := repo.NewProgressRepo(db)
	courseClient := client.NewcourseClient(
		cfg.CourseCfg.URL)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, progressRepo, courseClient, publisher)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentService)
	tokenManager := token.NewJWTManager(cfg.JWT.AccessSecret, cfg.JWT.Issuer)
	jwtMiddleware := middleware.JWT(tokenManager)
	routes.RegisterEnrollmentRoutes(app, enrollmentHandler, jwtMiddleware)
	log.Fatal(app.Listen(":" + cfg.App.Port))
}
