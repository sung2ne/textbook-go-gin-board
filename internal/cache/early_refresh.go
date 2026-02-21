package cache

import (
	"context"
	"time"
)

type CachedItem struct {
    Data      interface{} `json:"data"`
    ExpiresAt time.Time   `json:"expires_at"`
    RefreshAt time.Time   `json:"refresh_at"`  // 이 시간 이후 백그라운드 갱신
}

func GetWithEarlyRefresh[T any](
    ctx context.Context,
    key string,
    ttl time.Duration,
    refreshBefore time.Duration,  // 만료 전 이 시간부터 갱신 시작
    fetch func() (T, error),
) (T, error) {
    var item CachedItem
    var result T

    err := GetRedis(ctx, key, &item)
    if err == nil {
        result = item.Data.(T)

        // 갱신 시간이 지났으면 백그라운드에서 갱신
        if time.Now().After(item.RefreshAt) {
            go func() {
                newData, err := fetch()
                if err == nil {
                    newItem := CachedItem{
                        Data:      newData,
                        ExpiresAt: time.Now().Add(ttl),
                        RefreshAt: time.Now().Add(ttl - refreshBefore),
                    }
                    SetRedis(context.Background(), key, newItem, ttl)
                }
            }()
        }

        return result, nil
    }

    // 캐시 미스
    result, err = fetch()
    if err != nil {
        return result, err
    }

    item = CachedItem{
        Data:      result,
        ExpiresAt: time.Now().Add(ttl),
        RefreshAt: time.Now().Add(ttl - refreshBefore),
    }
    SetRedis(ctx, key, item, ttl)

    return result, nil
}
