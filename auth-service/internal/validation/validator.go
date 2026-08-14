package validation

import (
	"errors"
	"fmt"

	PV "github.com/go-playground/validator/v10"
)

var validate = PV.New(
	PV.WithRequiredStructEnabled(),
)

type FieldError struct {
	Field   string
	Code    string
	Message string
}

type Error struct {
	Fields []FieldError
}

func (e *Error) Error() string {
	return "request validation failed"
}

func ValidateStruct(value any) error {
	err := validate.Struct(value)
	if err == nil {
		return nil
	}

	var validationErrors PV.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return fmt.Errorf("validate request: %w", err)
	}

	fields := make([]FieldError, 0, len(validationErrors))

	for _, fieldErr := range validationErrors {
		fields = append(fields, FieldError{
			Field:   fieldErr.Field(),
			Code:    fieldErr.Tag(),
			Message: validationMessage(fieldErr),
		})
	}

	return &Error{
		Fields: fields,
	}
}

func validationMessage(err PV.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"

	case "email":
		return err.Field() + " must be a valid email"

	case "min":
		return err.Field() + " "
	case "max":
		return err.Field() + " "

	default:
		return err.Field() + " has an invalid value"
	}
}
