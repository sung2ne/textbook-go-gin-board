package auth

import (
    "context"

    "goboardapi/internal/domain"
    "goboardapi/internal/middleware"
)

// Checker 권한 체크 유틸리티
type Checker struct{}

func NewChecker() *Checker {
    return &Checker{}
}

// IsAuthenticated 인증 여부 확인
func (c *Checker) IsAuthenticated(ctx context.Context) bool {
    _, ok := middleware.GetUserFromContext(ctx)
    return ok
}

// GetCurrentUserID 현재 사용자 ID 반환
func (c *Checker) GetCurrentUserID(ctx context.Context) (uint, bool) {
    claims, ok := middleware.GetUserFromContext(ctx)
    if !ok {
        return 0, false
    }
    return claims.UserID, true
}

// GetCurrentRole 현재 사용자 역할 반환
func (c *Checker) GetCurrentRole(ctx context.Context) (domain.Role, bool) {
    claims, ok := middleware.GetUserFromContext(ctx)
    if !ok {
        return "", false
    }
    return domain.Role(claims.Role), true
}

// IsOwner 리소스 소유자인지 확인
func (c *Checker) IsOwner(ctx context.Context, ownerID uint) bool {
    userID, ok := c.GetCurrentUserID(ctx)
    if !ok {
        return false
    }
    return userID == ownerID
}

// IsAdmin 관리자인지 확인
func (c *Checker) IsAdmin(ctx context.Context) bool {
    role, ok := c.GetCurrentRole(ctx)
    if !ok {
        return false
    }
    return role == domain.RoleAdmin
}

// CanModify 수정 권한 확인 (본인 또는 관리자)
func (c *Checker) CanModify(ctx context.Context, ownerID uint) bool {
    return c.IsOwner(ctx, ownerID) || c.IsAdmin(ctx)
}

// HasRole 특정 역할 이상인지 확인
func (c *Checker) HasRole(ctx context.Context, requiredRole domain.Role) bool {
    role, ok := c.GetCurrentRole(ctx)
    if !ok {
        return false
    }
    return role.HasPermission(requiredRole)
}

// HasPermission 특정 권한이 있는지 확인
func (c *Checker) HasPermission(ctx context.Context, permission domain.Permission) bool {
    role, ok := c.GetCurrentRole(ctx)
    if !ok {
        return false
    }
    return domain.HasPermission(role, permission)
}
