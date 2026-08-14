package dto

type ErrorResponse struct {
	Error ErrorDetails `json:"course_errors"`
}
type ErrorDetails struct {
	Code      string             `json:"code"`
	Message   string             `json:"message"`
	RequestID string             `json:"request_id,omitempty"`
	Fields    []FieldErrorInside `json:"fields,omitempty"`
}
type FieldErrorInside struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
