package domain

// RoleGuest 게스트 역할
const RoleGuest Role = "guest"

// Level 역할 레벨 반환
func (r Role) Level() int {
	switch r {
	case RoleAdmin:
		return 100
	case RoleUser:
		return 10
	case RoleGuest:
		return 0
	default:
		return 0
	}
}

// HasPermission 권한 확인
func (r Role) HasPermission(permission Permission) bool {
	permissions, ok := RolePermissions[r]
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

// ParseRole 문자열에서 역할 파싱
func ParseRole(s string) Role {
	switch s {
	case string(RoleAdmin):
		return RoleAdmin
	case string(RoleUser):
		return RoleUser
	default:
		return RoleGuest
	}
}
