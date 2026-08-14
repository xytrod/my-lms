package main

import (
	"log"
	"main/auth-service/internal/client"
	"main/auth-service/internal/config"
	"main/auth-service/internal/handler"
	"main/auth-service/internal/hashing"
	"main/auth-service/internal/infra/database"
	"main/auth-service/internal/middleware"
	"main/auth-service/internal/repo"
	"main/auth-service/internal/routes"
	"main/auth-service/internal/service"
	"main/auth-service/internal/token"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal("database connection", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatal("database migration", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("sql database", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Println("sql database close", err)
		}
	}()

	authRepo := repo.NewAuthRepository(db)
	sessionRepo := repo.NewSessionRepository(db)
	passwordHashed := hashing.NewBcryptHash(bcrypt.DefaultCost)
	tokenManager := token.NewJWTManager(cfg.JWT)
	userClient := client.NewHTTPClient(cfg.UserService.BaseURL, cfg.UserService.APIKey, cfg.UserService.Timeout)
	authService := service.NewAuthService(authRepo, sessionRepo, userClient, passwordHashed, tokenManager)
	authHandler := handler.NewAuthHandler(authService)
	app := fiber.New(fiber.Config{
		ErrorHandler: handler.ErrorHandler,
	})
	app.Use(recover.New())
	app.Use(middleware.RequestID)
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"health": "ok",
		})
	})
	routes.RegisterAuthRoutes(app, authHandler, middleware.JWT(tokenManager))
	log.Fatal(app.Listen(":" + cfg.App.Port))

}
