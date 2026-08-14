package main

import (
	"log"
	"main/api-gateway/internal/config"
	"main/api-gateway/internal/proxy"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173",
			"http://localhost:3000",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "X-Request-ID", "Content-Type", "Authorization", "Accept"},
	}))
	app.Use(logger.New())
	p := proxy.New()
	app.Use("/auth", p.Forward(cfg.AuthURL))
	app.Use("/users", p.Forward(cfg.UserURL))
	app.Use("/courses", p.Forward(cfg.CourseURL))
	app.Use("/enrollments", p.Forward(cfg.EnrollmentURL))

	log.Fatal(app.Listen(":" + cfg.Port))
}
