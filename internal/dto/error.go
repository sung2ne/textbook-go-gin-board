
// AuthError 인증/인가 에러 응답
type AuthError struct {
    Error              string `json:"error"`
    Code               string `json:"code"`
    RequiredRole       string `json:"required_role,omitempty"`
    RequiredPermission string `json:"required_permission,omitempty"`
}

// UnauthorizedError 401 응답
func UnauthorizedError(message string) AuthError {
    return AuthError{
        Error: message,
        Code:  "UNAUTHORIZED",
    }
}

// ForbiddenError 403 응답
func ForbiddenError(message string, requiredRole domain.Role) AuthError {
    return AuthError{
        Error:        message,
        Code:         "FORBIDDEN",
        RequiredRole: string(requiredRole),
    }
}
