package domain

// Permission 권한
type Permission string

const (
    // 게시글 권한
    PermissionPostCreate Permission = "post:create"
    PermissionPostRead   Permission = "post:read"
    PermissionPostUpdate Permission = "post:update"
    PermissionPostDelete Permission = "post:delete"
    PermissionPostManage Permission = "post:manage" // 모든 게시글 관리

    // 댓글 권한
    PermissionCommentCreate Permission = "comment:create"
    PermissionCommentRead   Permission = "comment:read"
    PermissionCommentUpdate Permission = "comment:update"
    PermissionCommentDelete Permission = "comment:delete"
    PermissionCommentManage Permission = "comment:manage"

    // 사용자 권한
    PermissionUserRead   Permission = "user:read"
    PermissionUserManage Permission = "user:manage"
)

// RolePermissions 역할별 권한 매핑
var RolePermissions = map[Role][]Permission{
    RoleGuest: {
        PermissionPostRead,
        PermissionCommentRead,
    },
    RoleUser: {
        PermissionPostRead,
        PermissionPostCreate,
        PermissionPostUpdate,
        PermissionPostDelete,
        PermissionCommentRead,
        PermissionCommentCreate,
        PermissionCommentUpdate,
        PermissionCommentDelete,
    },
    RoleAdmin: {
        PermissionPostRead,
        PermissionPostCreate,
        PermissionPostUpdate,
        PermissionPostDelete,
        PermissionPostManage,
        PermissionCommentRead,
        PermissionCommentCreate,
        PermissionCommentUpdate,
        PermissionCommentDelete,
        PermissionCommentManage,
        PermissionUserRead,
        PermissionUserManage,
    },
}

// HasPermission 역할이 특정 권한을 가지고 있는지 확인
func HasPermission(role Role, permission Permission) bool {
    permissions, ok := RolePermissions[role]
    if !ok {
        return false
    }

    for _, p := range permissions {
        if p == permission {
            return true
        }
    }

    return false
}
