package validation

import (
	"errors"
	"fmt"
	"strings"

	PV "github.com/go-playground/validator/v10"
)

var Validate = PV.New(PV.WithRequiredStructEnabled())

func validationCode(tag string) string {
	switch tag {
	case "required":
		return "required"
	case "min":
		return "min"
	case "max":
		return "max"
	case "email":
		return "email"
	default:
		return "invalid"
	}
}
func validationMessage(field string, err PV.FieldError) string {
	switch field {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be greater than 7"
	case "max":
		return field + " must be less than or equal to 25"
	default:
		return field + " has an invalid value"
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
func ValidateStruct(value any) error {
	err := Validate.Struct(value)
	if err == nil {
		return nil
	}
	var validErrors PV.ValidationErrors
	if !errors.As(err, &validErrors) {
		return fmt.Errorf("validate request: %w", err)
	}
	fields := make([]FieldError, 0, len(validErrors))
	for _, fieldErr := range validErrors {
		field := toSnakeCase(fieldErr.Field())
		fields = append(fields, FieldError{
			Field:   field,
			Code:    validationCode(fieldErr.Tag()),
			Message: validationMessage(field, fieldErr),
		})
	}
	return &Error{
		Fields: fields,
	}
}
