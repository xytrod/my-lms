package handler

import (
	"main/enrollment_progress-service/internal/middleware"
	"main/enrollment_progress-service/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type EnrollmentHandler struct {
	service service.EnrollmentService
}

func NewEnrollmentHandler(service service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{service: service}
}
func getUserID(c fiber.Ctx) (uuid.UUID, error) {
	userID := fiber.Locals[uuid.UUID](c, middleware.UserIDKey)
	if userID == uuid.Nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "user id missing")
	}
	return userID, nil
}
func (h *EnrollmentHandler) Enroll(c fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}
	courseID, err := uuid.Parse(c.Params("courseID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid course id")
	}
	enrollment, err := h.service.Enroll(c.Context(), userID, courseID)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(enrollment)
}
func (h *EnrollmentHandler) CompletedLesson(c fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}
	courseID, err := uuid.Parse(c.Params("courseID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid course id")
	}
	lessonID, err := uuid.Parse(c.Params("lessonID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid lesson_id")
	}
	if err := h.service.CompletedLesson(c.Context(), userID, courseID, lessonID); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
func (h *EnrollmentHandler) GetProgress(c fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}
	courseID, err := uuid.Parse(c.Params("courseID"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid course id")
	}
	progress, err := h.service.GetProgress(c.Context(), userID, courseID)
	if err != nil {
		return err
	}
	return c.JSON(progress)
}
func (h *EnrollmentHandler) ListMy(c fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		return err
	}
	limit := fiber.Query[int](c, "limit", 10)
	offset := fiber.Query[int](c, "offset", 0)
	enrollments, err := h.service.ListMy(c.Context(), userID, limit, offset)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(enrollments)
}
