package cache

import (
	"sync/atomic"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var memoryCache *gocache.Cache

var (
	hits   int64
	misses int64
)

// InitMemoryCache - 인메모리 캐시 초기화
func InitMemoryCache() {
	memoryCache = gocache.New(5*time.Minute, 10*time.Minute)
}

// Set - 캐시에 저장
func Set(key string, value interface{}, expiration time.Duration) {
	memoryCache.Set(key, value, expiration)
}

// Get - 캐시에서 조회
func Get(key string) (interface{}, bool) {
	return memoryCache.Get(key)
}

// Delete - 캐시에서 삭제
func Delete(key string) {
	memoryCache.Delete(key)
}

// Flush - 전체 캐시 삭제
func Flush() {
	memoryCache.Flush()
}

// GetWithStats - 통계 포함 조회
func GetWithStats(key string) (interface{}, bool) {
	value, found := memoryCache.Get(key)
	if found {
		atomic.AddInt64(&hits, 1)
	} else {
		atomic.AddInt64(&misses, 1)
	}
	return value, found
}

// GetStats - 캐시 통계
func GetStats() (hitCount, missCount int64, hitRate float64) {
	h := atomic.LoadInt64(&hits)
	m := atomic.LoadInt64(&misses)
	total := h + m
	if total > 0 {
		hitRate = float64(h) / float64(total) * 100
	}
	return h, m, hitRate
}

// ResetStats - 통계 초기화
func ResetStats() {
	atomic.StoreInt64(&hits, 0)
	atomic.StoreInt64(&misses, 0)
}
