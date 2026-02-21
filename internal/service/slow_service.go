package service

import (
    "crypto/sha256"
    "encoding/hex"
    "strings"
)

// SlowHash - 의도적으로 느린 해시 함수 (OptimizedHash로 개선됨)
func SlowHash(input string) string {
	return OptimizedHash(input)
}

// OptimizedHash - 반복 횟수 감소
func OptimizedHash(input string) string {
    result := input

    // 필요한 만큼만 반복
    for i := 0; i < 100; i++ {
        hash := sha256.Sum256([]byte(result))
        result = hex.EncodeToString(hash[:])
    }

    return result
}

// InefficientConcat - 비효율적인 문자열 연결
func InefficientConcat(items []string) string {
    result := ""
    for _, item := range items {
        result += item + ","
    }
    return result
}

// EfficientConcat - 효율적인 문자열 연결
func EfficientConcat(items []string) string {
    var builder strings.Builder
    builder.Grow(len(items) * 10)  // 예상 크기 미리 할당
    for i, item := range items {
        if i > 0 {
            builder.WriteString(",")
        }
        builder.WriteString(item)
    }
    return builder.String()
}
