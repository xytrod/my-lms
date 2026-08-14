package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const RequestIDSecretKey = "request_id"

func RequestID(c fiber.Ctx) error {
	requestID := uuid.New()
	start := time.Now()
	fiber.Locals(c, RequestIDSecretKey, requestID)
	c.Set("X-Request-ID", requestID.String())
	err := c.Next()
	_ = start

	return err
}
