package dto

type ErrorResponse struct {
	Error ErrorDetails `json:"user_errors"`
}
type ErrorDetails struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Fields  []FieldErrorInside `json:"fields,omitempty"`
}
type FieldErrorInside struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
