package auth

import (
	"context"

	"goboardapi/internal/domain"
)

// Checker 권한 체크 유틸
type Checker struct{}

// NewChecker 체커 생성
func NewChecker() *Checker {
	return &Checker{}
}

// IsAuthenticated 인증 여부 확인
func (c *Checker) IsAuthenticated(ctx context.Context) bool {
	_, ok := ClaimsFromContext(ctx)
	return ok
}

// IsOwner 리소스 소유자 확인
func (c *Checker) IsOwner(ctx context.Context, ownerID uint) bool {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return false
	}
	return claims.UserID == ownerID
}

// IsAdmin 관리자 확인
func (c *Checker) IsAdmin(ctx context.Context) bool {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return false
	}
	return claims.Role == string(domain.RoleAdmin)
}

// CanModify 수정 권한 확인 (소유자 또는 관리자)
func (c *Checker) CanModify(ctx context.Context, ownerID uint) bool {
	return c.IsOwner(ctx, ownerID) || c.IsAdmin(ctx)
}

// HasPermission 권한 확인
func (c *Checker) HasPermission(ctx context.Context, permission domain.Permission) bool {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return false
	}
	role := domain.ParseRole(claims.Role)
	return role.HasPermission(permission)
}

// ClaimsFromContext 컨텍스트에서 클레임 추출 (middleware 패키지 순환 참조 방지)
func ClaimsFromContext(ctx context.Context) (*CustomClaims, bool) {
	claims, ok := ctx.Value("user").(*CustomClaims)
	return claims, ok
}
