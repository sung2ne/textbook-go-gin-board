package apperror

import (
	"fmt"
)

type AppError struct {
	HTTPStatus int
	Code       string
	Message    string
	Detail     string
	Err        error
	Fields     map[string]string
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
