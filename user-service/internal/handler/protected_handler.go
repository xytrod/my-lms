package handler

import "github.com/gofiber/fiber/v3"

func Protected(c fiber.Ctx) error {
	protected := fiber.Locals[bool](c, "client_authenticated")
	if !protected {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}
	return c.JSON(fiber.Map{
		"message":       "protected resource",
		"authenticated": true,
	})
}
