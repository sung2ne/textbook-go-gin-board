package auth

import "errors"

var (
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrNoPermission     = errors.New("no permission")
)

// RequireAuth 인증 필수 확인
func RequireAuth(claims *CustomClaims) error {
	if claims == nil {
		return ErrNotAuthenticated
	}
	return nil
}

// RequireOwner 소유자 필수 확인
func RequireOwner(claims *CustomClaims, ownerID uint) error {
	if claims == nil {
		return ErrNotAuthenticated
	}
	if claims.UserID != ownerID {
		return ErrNoPermission
	}
	return nil
}

// RequireAdmin 관리자 필수 확인
func RequireAdmin(claims *CustomClaims) error {
	if claims == nil {
		return ErrNotAuthenticated
	}
	if claims.Role != "admin" {
		return ErrNoPermission
	}
	return nil
}
