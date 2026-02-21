package cache

import (
    "time"
)

// TypedCache - 제네릭 캐시
type TypedCache[T any] struct {
    prefix string
}

// NewTypedCache - 타입 캐시 생성
func NewTypedCache[T any](prefix string) *TypedCache[T] {
    return &TypedCache[T]{prefix: prefix}
}

// Get - 타입 안전 조회
func (c *TypedCache[T]) Get(key string) (T, bool) {
    fullKey := c.prefix + ":" + key
    var zero T

    value, found := memoryCache.Get(fullKey)
    if !found {
        return zero, false
    }

    typed, ok := value.(T)
    if !ok {
        return zero, false
    }

    return typed, true
}

// Set - 타입 안전 저장
func (c *TypedCache[T]) Set(key string, value T, expiration time.Duration) {
    fullKey := c.prefix + ":" + key
    memoryCache.Set(fullKey, value, expiration)
}

// Delete - 삭제
func (c *TypedCache[T]) Delete(key string) {
    fullKey := c.prefix + ":" + key
    memoryCache.Delete(fullKey)
}
