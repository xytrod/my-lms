package validation

import (
	"errors"
	"fmt"
	"strings"

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
		field := toSnakeCase(fieldErr.Field())
		fields = append(fields, FieldError{
			Field:   field,
			Code:    fieldErr.Tag(),
			Message: fieldErr.Error(),
		})
	}

	return &Error{
		Fields: fields,
	}
}
func toSnakeCase(value string) string {
	var res strings.Builder
	for i, charact := range value {
		if charact >= 'A' && charact <= 'Z' {
			if i > 0 {
				res.WriteRune('_')
			}

			res.WriteRune(charact + ('a' - 'A'))
			continue
		}
		res.WriteRune(charact)
	}
	return res.String()

}

func validationMessage(field string, err PV.FieldError) string {
	switch err.Tag() {
	case "required":
		return field + " is required"
	case "min":
		return err.Field() + " must have at least 20 characters"
	case "max":
		return err.Field() + "must have less than 60 characters"
	case "oneof":
		return field + " invalid level"
	default:
		return field + " has an invalid value"
	}
}
