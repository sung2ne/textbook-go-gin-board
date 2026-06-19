package retry

import (
    "context"
    "fmt"
    "time"
)

type Config struct {
    MaxRetries int
    Delay      time.Duration
    Timeout    time.Duration
}

func WithRetry[T any](ctx context.Context, cfg Config, fn func(context.Context) (T, error)) (T, error) {
    var result T
    var lastErr error

    for i := 0; i <= cfg.MaxRetries; i++ {
        // 각 시도에 타임아웃 적용
        attemptCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)

        result, lastErr = fn(attemptCtx)
        cancel()

        if lastErr == nil {
            return result, nil
        }

        if i < cfg.MaxRetries {
            fmt.Printf("시도 %d 실패, %v 후 재시도: %v\n", i+1, cfg.Delay, lastErr)

            select {
            case <-ctx.Done():
                return result, ctx.Err()
            case <-time.After(cfg.Delay):
            }
        }
    }

    return result, fmt.Errorf("최대 재시도 횟수 초과: %w", lastErr)
}
