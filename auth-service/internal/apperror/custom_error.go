package apperror

import "fmt"

type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Err)
	}
	return e.Code
}
func (e *AppError) Unwrap() error {
	return e.Err
}
