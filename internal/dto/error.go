package dto

// AuthError 인증 에러 응답
type AuthError struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// UnauthorizedError 미인증 에러
func UnauthorizedError(message string) *AuthError {
	return &AuthError{
		Error: message,
		Code:  "UNAUTHORIZED",
	}
}

// ForbiddenError 권한 없음 에러
func ForbiddenError(message string) *AuthError {
	return &AuthError{
		Error: message,
		Code:  "FORBIDDEN",
	}
}
