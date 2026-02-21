package cache

import (
    "context"
    "time"
)

// AcquireLock - 분산 락 획득
func AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
    return redisClient.SetNX(ctx, "lock:"+key, "1", ttl).Result()
}

// ReleaseLock - 분산 락 해제
func ReleaseLock(ctx context.Context, key string) error {
    return redisClient.Del(ctx, "lock:"+key).Err()
}

// WithLock - 락과 함께 실행
func WithLock(ctx context.Context, key string, ttl time.Duration, fn func() error) error {
    acquired, err := AcquireLock(ctx, key, ttl)
    if err != nil {
        return err
    }

    if !acquired {
        // 다른 프로세스가 이미 처리 중
        return nil
    }

    defer ReleaseLock(ctx, key)
    return fn()
}
