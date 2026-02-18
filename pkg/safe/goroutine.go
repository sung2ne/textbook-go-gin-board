package safe

import (
	"runtime/debug"

	"goboardapi/pkg/logger"
)

func Go(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				logger.Error().
					Interface("panic", err).
					Str("stack", stack).
					Msg("고루틴 패닉 복구")
			}
		}()
		fn()
	}()
}
