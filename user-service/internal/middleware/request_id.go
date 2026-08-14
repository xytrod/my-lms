package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const RequestIDSecretKey = "request_id"

func RequestIDMiddleware(c fiber.Ctx) error {
	requestID := uuid.New()
	start := time.Now()
	fiber.Locals(c, RequestIDSecretKey, requestID)
	c.Set("X-Request-ID", requestID.String())
	err := c.Next()
	duration := time.Since(start)
	log.Printf("Request finished: request_id=%s method=%s status=%s duration=%s",
		requestID, c.Method(), c.Response().StatusCode(), duration)

	return err
}
