package middleware

import (
	"main/course-service/internal/token"
	"strings"

	"github.com/gofiber/fiber/v3"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

func JWT(manager token.Manager) fiber.Handler {
	return func(c fiber.Ctx) error {
		authheader := c.Get("Authorization")
		if authheader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing Authorization header")
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(authheader, prefix) {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header")
		}
		raw := strings.TrimSpace(strings.TrimPrefix(authheader, prefix))
		if raw == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "access token is required")
		}
		claims, err := manager.ParseAccessToken(raw)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid, expired access token")
		}
		fiber.Locals(c, UserIDKey, claims.UserID)
		fiber.Locals(c, RoleKey, claims.Role)
		return c.Next()
	}
}
