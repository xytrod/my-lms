package handler

import (
	"errors"
	"log"
	"main/user-service/internal/dto"
	"main/user-service/internal/user_errors"
	"main/user-service/internal/validation"

	"github.com/gofiber/fiber/v3"
)

func Errorhandler(c fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	code := "internal_error"
	message := "internal server course_errors"
	switch {
	case errors.Is(err, user_errors.ErrUserNotFound):
		status = fiber.StatusNotFound
		code = "user_not_found"
		message = err.Error()
	case errors.Is(err, user_errors.ErrInvalidUserID):
		status = fiber.StatusBadRequest
		code = "invalid_user_id"
		message = err.Error()
	case errors.Is(err, user_errors.ErrEmailRequired):
		status = fiber.StatusBadRequest
		code = "email_required"
		message = err.Error()
	case errors.Is(err, user_errors.ErrInvalidEmail):
		status = fiber.StatusBadRequest
		code = "invalid_email"
		message = err.Error()
	case errors.Is(err, user_errors.ErrEmailAlreadyExists):
		status = fiber.StatusConflict
		code = "email_already_exists"
		message = err.Error()
	case errors.Is(err, user_errors.ErrUsernameRequired):
		status = fiber.StatusBadRequest
		code = "username_required"
		message = err.Error()
	case errors.Is(err, user_errors.ErrFirstNameRequired):
		status = fiber.StatusBadRequest
		code = "first_name_required"
		message = err.Error()
	case errors.Is(err, user_errors.ErrLastNameRequired):
		status = fiber.StatusBadRequest
		code = "last_name_required"
		message = err.Error()
	case errors.Is(err, user_errors.ErrUsernameAlreadyExists):
		status = fiber.StatusConflict
		code = "username_already_exists"
		message = err.Error()
	}
	var validationError *validation.Error
	if errors.As(err, &validationError) {
		fields := make([]dto.FieldErrorInside, len(validationError.Fields))
		for _, field := range validationError.Fields {
			fields = append(fields, dto.FieldErrorInside{
				Field:   field.Field,
				Code:    field.Code,
				Message: field.Message,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: dto.ErrorDetails{
				Code:    "validation_failed",
				Message: "validation failed",
				Fields:  fields,
			},
		})
	}
	var fiber_err *fiber.Error
	if errors.As(err, &fiber_err) {
		status = fiber_err.Code
		code = "http_error"
		message = fiber_err.Message
	}
	if status >= 500 {
		log.Println("internal course_errors: ", err)
	}
	return c.Status(status).JSON(dto.ErrorResponse{
		Error: dto.ErrorDetails{
			Code:    code,
			Message: message,
		},
	})
}
