package auth

import (
	"context"

	"goboardapi/internal/domain"
)

// PermissionBuilder 권한 빌더
type PermissionBuilder struct {
	ctx        context.Context
	conditions []func() bool
	mode       string // "any" or "all"
}

// NewPermissionBuilder 빌더 생성
func NewPermissionBuilder(ctx context.Context) *PermissionBuilder {
	return &PermissionBuilder{
		ctx:  ctx,
		mode: "all",
	}
}

// Any 조건 중 하나만 만족하면 됨
func (b *PermissionBuilder) Any() *PermissionBuilder {
	b.mode = "any"
	return b
}

// All 모든 조건을 만족해야 함
func (b *PermissionBuilder) All() *PermissionBuilder {
	b.mode = "all"
	return b
}

// IsOwner 소유자 조건 추가
func (b *PermissionBuilder) IsOwner(ownerID uint) *PermissionBuilder {
	b.conditions = append(b.conditions, func() bool {
		return NewChecker().IsOwner(b.ctx, ownerID)
	})
	return b
}

// IsAdmin 관리자 조건 추가
func (b *PermissionBuilder) IsAdmin() *PermissionBuilder {
	b.conditions = append(b.conditions, func() bool {
		return NewChecker().IsAdmin(b.ctx)
	})
	return b
}

// HasPermission 권한 조건 추가
func (b *PermissionBuilder) HasPermission(permission domain.Permission) *PermissionBuilder {
	b.conditions = append(b.conditions, func() bool {
		return NewChecker().HasPermission(b.ctx, permission)
	})
	return b
}

// Check 권한 확인
func (b *PermissionBuilder) Check() bool {
	if len(b.conditions) == 0 {
		return true
	}

	if b.mode == "any" {
		for _, cond := range b.conditions {
			if cond() {
				return true
			}
		}
		return false
	}

	// all mode
	for _, cond := range b.conditions {
		if !cond() {
			return false
		}
	}
	return true
}

// Require 권한 필수 확인
func (b *PermissionBuilder) Require() error {
	if !NewChecker().IsAuthenticated(b.ctx) {
		return ErrNotAuthenticated
	}
	if !b.Check() {
		return ErrNoPermission
	}
	return nil
}
