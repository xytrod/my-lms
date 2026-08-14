package main

import (
	"log"
	"main/course-service/internal/broker"
	"main/course-service/internal/config"
	"main/course-service/internal/handler"
	"main/course-service/internal/infra/database"
	"main/course-service/internal/middleware"
	"main/course-service/internal/repo"
	"main/course-service/internal/routes"
	"main/course-service/internal/service"
	"main/course-service/internal/token"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("load", err)
	}
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
	rabbit, err := broker.NewRabbitMQ(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("rabbit", err)
	}
	defer func() {
		if err := rabbit.Close(); err != nil {
			log.Fatal("rabbit close", err)
		}
	}()
	if err := rabbit.StartUserEnrolledConsumer(); err != nil {
		log.Fatal("start user enrolled consumer", err)
	}
	courseRepo := repo.NewCourseRepository(db)
	lessonRepo := repo.NewLessonRepository(db)
	courseService := service.NewCourseService(courseRepo, lessonRepo)
	courseHandler := handler.NewCourseHandler(courseService)
	tokenManager := token.NewJWTManager(string(cfg.JWT.AccessSecret), cfg.JWT.Issuer)
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.Errorhandler,
	})
	app.Use(recover.New())
	app.Use(middleware.RequestID)
	routes.RegisterCourseRoutes(app, courseHandler, middleware.JWT(tokenManager))
	for _, route := range app.GetRoutes(true) {
		log.Printf(
			"ROUTE: %-7s %s",
			route.Method,
			route.Path,
		)
	}
	log.Fatal(app.Listen(":" + cfg.App.Port))
}
