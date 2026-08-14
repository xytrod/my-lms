package routes

import (
	"main/auth-service/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(app *fiber.App, h *handler.AuthHandler, jwtMiddleware fiber.Handler) {
	auth := app.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)
	protected := auth.Group("")
	protected.Use(jwtMiddleware)
	protected.Get("/me", h.MyProfile)
}
