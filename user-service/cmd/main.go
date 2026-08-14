package main

import (
	"log"
	"main/user-service/internal/config"
	"main/user-service/internal/handler"
	"main/user-service/internal/infra/database"
	"main/user-service/internal/middleware"
	"main/user-service/internal/repository"
	"main/user-service/internal/routes"
	"main/user-service/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg := config.Load()
	log.Printf(
		"USER DB CONFIG host=%q port=%q db=%q user=%q",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.Username,
	)
	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal("migration failed", err)
	}
	sqldb, err := db.DB()
	if err != nil {
		log.Fatal("sql database failed", err)
	}
	defer func() {
		if err := sqldb.Close(); err != nil {
			log.Fatal("sql close failed", err)
		}
	}()
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.Errorhandler,
	})
	app.Use(recover.New())
	app.Use(middleware.RequestIDMiddleware)
	routes.RegisterUserRoutes(app, userHandler)
	protected := app.Group("/api")
	protected.Use(middleware.APIKey(cfg.App.APIKey))
	routes.RegisterProtectedRoutes(protected, userHandler)
	for _, route := range app.GetRoutes(true) {
		log.Printf(
			"ROUTE: %-7s %s",
			route.Method,
			route.Path,
		)
	}
	log.Fatal(app.Listen(":8080"))
}
