package cache

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

// CacheAside - 캐시 어사이드 패턴
func CacheAside[T any](
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

    // 2. 캐시 미스 - 데이터 가져오기
    result, err = fetch()
    if err != nil {
        return result, err
    }

    // 3. 캐시에 저장
    SetRedis(ctx, key, result, ttl)

    return result, nil
}
