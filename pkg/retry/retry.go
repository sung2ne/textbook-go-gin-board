package retry

import (
	"context"
	"fmt"
	"time"
)

// Config 재시도 설정
type Config struct {
	MaxRetries int
	Delay      time.Duration
	Timeout    time.Duration
}

// WithRetry 재시도 로직
func WithRetry[T any](ctx context.Context, cfg Config, fn func(context.Context) (T, error)) (T, error) {
	var result T
	var lastErr error

	for i := 0; i <= cfg.MaxRetries; i++ {
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
