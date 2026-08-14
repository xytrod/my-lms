package routes

import (
	"main/user-service/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterProtectedRoutes(group fiber.Router, h *handler.UserHandler) {
	users := group.Group("/users")
	users.Get("/", h.List)
	users.Get("/:id", h.GetByID)
	users.Patch("/:id", h.Update)
	users.Delete("/:id", h.Delete)
	group.Get("/protected", handler.Protected)
}
