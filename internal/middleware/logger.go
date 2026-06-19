package middleware

import (
    "log/slog"
    "time"

    "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        // 요청 처리
        c.Next()

        // 응답 후 로깅
        latency := time.Since(start)
        status := c.Writer.Status()

        attrs := []any{
            "status", status,
            "method", c.Request.Method,
            "path", path,
            "latency", latency.String(),
            "ip", c.ClientIP(),
        }

        if query != "" {
            attrs = append(attrs, "query", query)
        }

        if len(c.Errors) > 0 {
            attrs = append(attrs, "error", c.Errors.String())
        }

        switch {
        case status >= 500:
            slog.Error("요청 처리 실패", attrs...)
        case status >= 400:
            slog.Warn("클라이언트 오류", attrs...)
        default:
            slog.Info("요청 처리 완료", attrs...)
        }
    }
}
