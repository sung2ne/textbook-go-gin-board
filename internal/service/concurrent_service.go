package service

import (
    "sync"
    "time"
)

type Counter struct {
    mu    sync.RWMutex
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()

    time.Sleep(time.Millisecond)
    c.value++
}

// GetValue - 읽기 락 사용
func (c *Counter) GetValue() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.value
}
