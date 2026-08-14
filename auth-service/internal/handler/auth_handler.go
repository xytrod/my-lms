package handler

import (
	"main/auth-service/internal/dto"
	"main/auth-service/internal/middleware"
	"main/auth-service/internal/service"
	"main/auth-service/internal/validation"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	resp, err := h.service.Register(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	resp, err := h.service.Login(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req dto.RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	resp, err := h.service.Refresh(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req dto.LogoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}
	if err := h.service.Logout(c.Context(), req); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *AuthHandler) MyProfile(c fiber.Ctx) error {
	userID := fiber.Locals[uuid.UUID](c, middleware.UserIDKey)
	role := fiber.Locals[string](c, middleware.RoleKey)
	return c.JSON(fiber.Map{
		"user_id": userID,
		"role":    role,
	})
}
