package apperror

import "net/http"

func Wrap(err error, code string, message string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func WrapWithStatus(err error, status int, code string, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func (e *AppError) WithDetail(detail string) *AppError {
	e.Detail = detail
	return e
}

func (e *AppError) WithFields(fields map[string]string) *AppError {
	e.Fields = fields
	return e
}
