package cache

import (
    "context"
    "time"
)

// GetLayered - L1(메모리) → L2(Redis) → DB
func GetLayered[T any](
    ctx context.Context,
    key string,
    l1TTL time.Duration,
    l2TTL time.Duration,
    fetch func() (T, error),
) (T, error) {
    var result T

    // L1: 인메모리 캐시
    if cached, found := Get(key); found {
        return cached.(T), nil
    }

    // L2: Redis
    err := GetRedis(ctx, key, &result)
    if err == nil {
        // L1에 복사
        Set(key, result, l1TTL)
        return result, nil
    }

    // DB
    result, err = fetch()
    if err != nil {
        return result, err
    }

    // 양쪽에 저장
    Set(key, result, l1TTL)
    SetRedis(ctx, key, result, l2TTL)

    return result, nil
}
