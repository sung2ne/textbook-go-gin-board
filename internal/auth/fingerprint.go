package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

// GenerateFingerprint는 클라이언트 고유 식별자를 생성합니다.
func GenerateFingerprint(r *http.Request) string {
	userAgent := r.UserAgent()
	data := userAgent
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}
