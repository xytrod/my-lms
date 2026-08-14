package routes

import (
	"main/user-service/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(app *fiber.App, h *handler.UserHandler) {
	users := app.Group("/users")
	users.Post("/", h.Create)
	users.Get("/health", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "OK",
		})
	})
}
