package middleware

import (
	"github.com/gofiber/fiber/v3"
)

const expected_key = "secret"

func APIKey(expectedKey string) fiber.Handler {
	return func(c fiber.Ctx) error {
		apiKey := c.Get("X-API-KEY")
		if apiKey == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}
		if apiKey != expected_key {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid api key")
		}
		fiber.Locals(c, "client_authenticated", true)
		return c.Next()
	}

}
