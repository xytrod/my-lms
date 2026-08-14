package validation

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
type Error struct {
	Fields []FieldError
}

func (e *Error) Error() string {
	return "entry validation failed"
}
