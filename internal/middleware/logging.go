package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"goboardapi/pkg/logger"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		requestID := c.GetString("requestID")

		requestLogger := logger.Logger.With().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("client_ip", c.ClientIP()).
			Logger()

		c.Set("logger", requestLogger)

		requestLogger.Info().
			Str("user_agent", c.Request.UserAgent()).
			Str("query", query).
			Msg("요청 시작")

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		var event *zerolog.Event
		switch {
		case statusCode >= 500:
			event = requestLogger.Error()
		case statusCode >= 400:
			event = requestLogger.Warn()
		default:
			event = requestLogger.Info()
		}

		event.
			Int("status", statusCode).
			Dur("latency", latency).
			Int("body_size", c.Writer.Size()).
			Msg("요청 완료")
	}
}
