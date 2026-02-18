package middleware

import (
	"context"

	"goboardapi/internal/auth"

	"github.com/gin-gonic/gin"
)

type contextKey string

const userContextKey contextKey = "user"

// SetUserToContext는 사용자 정보를 context.Context에 저장합니다.
func SetUserToContext(ctx context.Context, claims *auth.CustomClaims) context.Context {
	return context.WithValue(ctx, userContextKey, claims)
}

// GetUserFromContext는 context.Context에서 사용자 정보를 추출합니다.
func GetUserFromContext(ctx context.Context) (*auth.CustomClaims, bool) {
	claims, ok := ctx.Value(userContextKey).(*auth.CustomClaims)
	return claims, ok
}

// GetCurrentUser는 Gin 컨텍스트에서 현재 사용자 정보를 추출합니다.
func GetCurrentUser(c *gin.Context) (*auth.CustomClaims, bool) {
	value, exists := c.Get(ContextUserKey)
	if !exists {
		return nil, false
	}

	claims, ok := value.(*auth.CustomClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}

// MustGetCurrentUser는 현재 사용자 정보를 추출합니다. 없으면 패닉입니다.
func MustGetCurrentUser(c *gin.Context) *auth.CustomClaims {
	claims, ok := GetCurrentUser(c)
	if !ok {
		panic("user not found in context - auth middleware not applied?")
	}
	return claims
}

// IsAuthenticated는 현재 요청이 인증되었는지 확인합니다.
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(ContextUserKey)
	return exists
}

// GetCurrentUserID는 현재 사용자 ID를 반환합니다.
func GetCurrentUserID(c *gin.Context) uint {
	claims, ok := GetCurrentUser(c)
	if !ok {
		return 0
	}
	return claims.UserID
}
