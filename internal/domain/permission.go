package domain

// Permission 권한 타입
type Permission string

const (
	PermissionPostCreate Permission = "post:create"
	PermissionPostUpdate Permission = "post:update"
	PermissionPostDelete Permission = "post:delete"
	PermissionPostRead   Permission = "post:read"

	PermissionCommentCreate Permission = "comment:create"
	PermissionCommentUpdate Permission = "comment:update"
	PermissionCommentDelete Permission = "comment:delete"

	PermissionUserManage Permission = "user:manage"
	PermissionUserDelete Permission = "user:delete"
)

// RolePermissions 역할별 권한 매핑
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermissionPostCreate, PermissionPostUpdate, PermissionPostDelete, PermissionPostRead,
		PermissionCommentCreate, PermissionCommentUpdate, PermissionCommentDelete,
		PermissionUserManage, PermissionUserDelete,
	},
	RoleUser: {
		PermissionPostCreate, PermissionPostUpdate, PermissionPostDelete, PermissionPostRead,
		PermissionCommentCreate, PermissionCommentUpdate, PermissionCommentDelete,
	},
	RoleGuest: {
		PermissionPostRead,
	},
}

// HasPermission 역할에 권한이 있는지 확인
func HasPermission(role Role, permission Permission) bool {
	return role.HasPermission(permission)
}
