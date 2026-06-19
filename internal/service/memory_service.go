package service

import (
    "sync"
    "time"
)

type cacheItem struct {
    data      []byte
    expiresAt time.Time
}

type BoundedCache struct {
    items    map[string]cacheItem
    mu       sync.RWMutex
    maxSize  int
}

func NewBoundedCache(maxSize int) *BoundedCache {
    c := &BoundedCache{
        items:   make(map[string]cacheItem),
        maxSize: maxSize,
    }

    // 주기적으로 만료 항목 정리
    go c.cleanup()

    return c
}

func (c *BoundedCache) Set(key string, data []byte, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // 크기 제한 확인
    if len(c.items) >= c.maxSize {
        c.evictOldest()
    }

    c.items[key] = cacheItem{
        data:      data,
        expiresAt: time.Now().Add(ttl),
    }
}

func (c *BoundedCache) Get(key string) ([]byte, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, ok := c.items[key]
    if !ok || time.Now().After(item.expiresAt) {
        return nil, false
    }

    return item.data, true
}

func (c *BoundedCache) evictOldest() {
    var oldestKey string
    var oldestTime time.Time

    for key, item := range c.items {
        if oldestKey == "" || item.expiresAt.Before(oldestTime) {
            oldestKey = key
            oldestTime = item.expiresAt
        }
    }

    if oldestKey != "" {
        delete(c.items, oldestKey)
    }
}

func (c *BoundedCache) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, item := range c.items {
            if now.After(item.expiresAt) {
                delete(c.items, key)
            }
        }
        c.mu.Unlock()
    }
}
