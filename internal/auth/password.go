package auth

import "golang.org/x/crypto/bcrypt"

// PasswordHasher 비밀번호 해싱 인터페이스
type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hashedPassword, password string) bool
}

type bcryptHasher struct {
    cost int
}

func NewBcryptHasher(cost int) PasswordHasher {
    return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
    return string(bytes), err
}

func (h *bcryptHasher) Compare(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
