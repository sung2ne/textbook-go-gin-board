package apperror

import (
	"fmt"
	"net/http"
)

func NotFound(resource string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s을(를) 찾을 수 없습니다", resource),
	}
}

func NotFoundWithID(resource string, id any) *AppError {
	return &AppError{
		HTTPStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    fmt.Sprintf("%s을(를) 찾을 수 없습니다", resource),
		Detail:     fmt.Sprintf("ID: %v", id),
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

func ValidationError(fields map[string]string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    "입력값이 올바르지 않습니다",
		Fields:     fields,
	}
}

func Unauthorized(message string) *AppError {
	if message == "" {
		message = "인증이 필요합니다"
	}
	return &AppError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

func Forbidden(message string) *AppError {
	if message == "" {
		message = "접근 권한이 없습니다"
	}
	return &AppError{
		HTTPStatus: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    message,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

func InternalError(err error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    "서버 내부 오류가 발생했습니다",
		Err:        err,
	}
}
