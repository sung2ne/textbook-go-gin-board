package cache

import (
    "sync/atomic"
)

var (
    hits   int64
    misses int64
)

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
