
// ParseRole 문자열을 Role로 변환
func ParseRole(s string) (Role, error) {
    role := Role(s)
    if !role.IsValid() {
        return "", fmt.Errorf("invalid role: %s", s)
    }
    return role, nil
}

// MustParseRole 문자열을 Role로 변환 (패닉 버전)
func MustParseRole(s string) Role {
    role, err := ParseRole(s)
    if err != nil {
        panic(err)
    }
    return role
}
