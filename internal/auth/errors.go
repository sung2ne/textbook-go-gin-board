package auth

import (
    "context"
    "errors"
)

var (
    ErrNotAuthenticated = errors.New("인증이 필요합니다")
    ErrNoPermission     = errors.New("권한이 없습니다")
    ErrNotOwner         = errors.New("소유자가 아닙니다")
)

// RequireAuth 인증 필수 확인
func RequireAuth(ctx context.Context) error {
    if !NewChecker().IsAuthenticated(ctx) {
        return ErrNotAuthenticated
    }
    return nil
}

// RequireOwner 소유자 필수 확인
func RequireOwner(ctx context.Context, ownerID uint) error {
    checker := NewChecker()
    if !checker.IsAuthenticated(ctx) {
        return ErrNotAuthenticated
    }
    if !checker.IsOwner(ctx, ownerID) {
        return ErrNotOwner
    }
    return nil
}

// RequireOwnerOrAdmin 소유자 또는 관리자 필수 확인
func RequireOwnerOrAdmin(ctx context.Context, ownerID uint) error {
    checker := NewChecker()
    if !checker.IsAuthenticated(ctx) {
        return ErrNotAuthenticated
    }
    if !checker.CanModify(ctx, ownerID) {
        return ErrNoPermission
    }
    return nil
}

// RequireAdmin 관리자 필수 확인
func RequireAdmin(ctx context.Context) error {
    checker := NewChecker()
    if !checker.IsAuthenticated(ctx) {
        return ErrNotAuthenticated
    }
    if !checker.IsAdmin(ctx) {
        return ErrNoPermission
    }
    return nil
}

// RequirePermission 특정 권한 필수 확인
func RequirePermission(ctx context.Context, permission domain.Permission) error {
    checker := NewChecker()
    if !checker.IsAuthenticated(ctx) {
        return ErrNotAuthenticated
    }
    if !checker.HasPermission(ctx, permission) {
        return ErrNoPermission
    }
    return nil
}
