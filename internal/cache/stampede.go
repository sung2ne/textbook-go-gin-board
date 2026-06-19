package cache

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

func GetWithLock[T any](
    ctx context.Context,
    key string,
    ttl time.Duration,
    fetch func() (T, error),
) (T, error) {
    var result T

    // 1. 캐시 조회
    err := GetRedis(ctx, key, &result)
    if err == nil {
        return result, nil
    }

    if err != redis.Nil {
        return result, err
    }

    // 2. 락 획득 시도
    lockKey := "lock:" + key
    acquired, _ := AcquireLock(ctx, lockKey, 10*time.Second)

    if !acquired {
        // 락 획득 실패 - 잠시 대기 후 재시도
        time.Sleep(100 * time.Millisecond)
        return GetWithLock(ctx, key, ttl, fetch)
    }
    defer ReleaseLock(ctx, lockKey)

    // 3. 다시 캐시 확인 (다른 프로세스가 갱신했을 수 있음)
    err = GetRedis(ctx, key, &result)
    if err == nil {
        return result, nil
    }

    // 4. DB에서 가져와서 캐시
    result, err = fetch()
    if err != nil {
        return result, err
    }

    SetRedis(ctx, key, result, ttl)
    return result, nil
}
