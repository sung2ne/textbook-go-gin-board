package middleware

import (
	"errors"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"syscall"

	"github.com/gin-gonic/gin"
	"goboardapi/pkg/logger"
	"goboardapi/pkg/problem"
)

type RecoveryConfig struct {
	EnableStackTrace bool
	NotifyFunc       func(err any, stack string)
}

func isBrokenPipe(err any) bool {
	if netErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := netErr.Err.(*os.SyscallError); ok {
			if errors.Is(sysErr.Err, syscall.EPIPE) ||
				errors.Is(sysErr.Err, syscall.ECONNRESET) {
				return true
			}
		}
	}
	return false
}

func Recovery(cfg RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if isBrokenPipe(err) {
					log := logger.FromGin(c)
					log.Warn().Interface("error", err).Msg("브로큰 파이프")
					c.Abort()
					return
				}

				stack := string(debug.Stack())

				log := logger.FromGin(c)
				event := log.Error().
					Interface("panic", err).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path)

				if cfg.EnableStackTrace {
					event = event.Str("stack", stack)
				}
				event.Msg("패닉 복구")

				if cfg.NotifyFunc != nil {
					go cfg.NotifyFunc(err, stack)
				}

				pd := problem.InternalError(c.Request.URL.Path)
				c.Header("Content-Type", problem.ContentType)
				c.AbortWithStatusJSON(http.StatusInternalServerError, pd)
			}
		}()

		c.Next()
	}
}

func DefaultRecovery() gin.HandlerFunc {
	return Recovery(RecoveryConfig{
		EnableStackTrace: true,
	})
}
