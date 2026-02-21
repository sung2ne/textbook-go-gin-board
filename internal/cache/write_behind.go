package cache

import (
    "context"
    "sync"
    "time"
)

type WriteBuffer struct {
    mu      sync.Mutex
    pending map[string]interface{}
    flush   func(ctx context.Context, key string, value interface{}) error
}

func NewWriteBuffer(flush func(ctx context.Context, key string, value interface{}) error) *WriteBuffer {
    wb := &WriteBuffer{
        pending: make(map[string]interface{}),
        flush:   flush,
    }

    // 주기적으로 DB에 반영
    go wb.backgroundFlush()

    return wb
}

func (wb *WriteBuffer) Write(key string, value interface{}) {
    wb.mu.Lock()
    defer wb.mu.Unlock()
    wb.pending[key] = value
}

func (wb *WriteBuffer) backgroundFlush() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        wb.mu.Lock()
        pending := wb.pending
        wb.pending = make(map[string]interface{})
        wb.mu.Unlock()

        ctx := context.Background()
        for key, value := range pending {
            wb.flush(ctx, key, value)
        }
    }
}
