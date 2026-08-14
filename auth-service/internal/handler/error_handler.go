package handler

import (
	"errors"
	"log"
	"main/auth-service/internal/apperror"
	"main/auth-service/internal/dto"
	"main/auth-service/internal/validation"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	requestID := getRequestID(c)
	var ValidErr *validation.Error

	if errors.As(err, &ValidErr) {
		fields := make([]dto.FieldErrorInside, 0, len(ValidErr.Fields))
		for _, field := range ValidErr.Fields {
			fields = append(fields, dto.FieldErrorInside{
				Field:   strings.ToLower(field.Field),
				Code:    field.Code,
				Message: field.Message,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Error: dto.ErrorDetails{
					Code:      "validation failed",
					Message:   "request validation failed",
					RequestID: getRequestID(c),
					Fields:    fields,
				},
			})
	}
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		if appErr.Status >= fiber.StatusInternalServerError {
			log.Printf(
				"request_id=%s code=%s internal_error=%v",
				requestID,
				appErr.Code,
				appErr.Err)
		}
		return c.Status(fiber.StatusBadRequest).JSON(
			dto.ErrorResponse{
				Error: dto.ErrorDetails{
					Code:      appErr.Code,
					Message:   appErr.Message,
					RequestID: requestID,
				},
			})
	}
	var fiber_err *fiber.Error
	if errors.As(err, &fiber_err) && fiber_err != nil {
		return c.Status(fiber_err.Code).JSON(
			dto.ErrorResponse{
				Error: dto.ErrorDetails{
					Code:      "http_error",
					Message:   fiber_err.Message,
					RequestID: requestID,
				},
			})
	}
	log.Println("unknown course_errors", err, requestID)
	return c.Status(fiber.StatusInternalServerError).JSON(
		dto.ErrorResponse{
			Error: dto.ErrorDetails{
				Code:      "internal_error",
				Message:   "internal server course_errors",
				RequestID: requestID,
			},
		})
}
func getRequestID(c fiber.Ctx) string {
	requestID := fiber.Locals[uuid.UUID](c, "request_id")
	if requestID == uuid.Nil {
		return ""
	}
	return requestID.String()
}
