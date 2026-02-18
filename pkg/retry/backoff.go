package retry

import (
	"context"
	"time"
)

// BackoffConfig 지수 백오프 설정
type BackoffConfig struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	MaxRetries   int
	Multiplier   float64
}

// WithExponentialBackoff 지수 백오프 재시도
func WithExponentialBackoff[T any](
	ctx context.Context,
	cfg BackoffConfig,
	fn func(context.Context) (T, error),
) (T, error) {
	var result T
	var lastErr error
	delay := cfg.InitialDelay

	for i := 0; i <= cfg.MaxRetries; i++ {
		result, lastErr = fn(ctx)
		if lastErr == nil {
			return result, nil
		}

		if i < cfg.MaxRetries {
			select {
			case <-ctx.Done():
				return result, ctx.Err()
			case <-time.After(delay):
			}

			delay = time.Duration(float64(delay) * cfg.Multiplier)
			if delay > cfg.MaxDelay {
				delay = cfg.MaxDelay
			}
		}
	}

	return result, lastErr
}
