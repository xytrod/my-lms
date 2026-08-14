package middleware

import (
	"main/auth-service/internal/token"
	"strings"

	"github.com/gofiber/fiber/v3"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "role"
)

func JWT(manager token.Manager) fiber.Handler {
	return func(c fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Authorization required")
		}
		prefix := "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid string")
		}
		raw := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
		if raw == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "access token is required")
		}
		claims, err := manager.ParseAccessToken(raw)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token (or expired)")
		}
		fiber.Locals(c, UserIDKey, claims.UserID)
		fiber.Locals(c, RoleKey, claims.Role)
		return c.Next()
	}
}
