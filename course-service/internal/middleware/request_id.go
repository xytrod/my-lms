package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequestID(c fiber.Ctx) error {
	reqID := uuid.New()
	start := time.Now()
	fiber.Locals(c, "request_id", reqID)
	c.Set("X-Request-ID", reqID.String())
	err := c.Next()
	log.Printf(
		"request_id=%s method=%s path=%s status=%d duration=%s",
		reqID,
		c.Method(),
		c.Path(),
		c.Response().StatusCode(),
		time.Since(start),
	)
	return err
}
