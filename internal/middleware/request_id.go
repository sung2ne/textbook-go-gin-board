package middleware

import (
    "log/slog"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type ctxKey string

const RequestIDKey ctxKey = "request_id"

func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }

        // Gin 컨텍스트에 저장
        c.Set(string(RequestIDKey), requestID)

        // 응답 헤더에도 추가
        c.Header("X-Request-ID", requestID)

        c.Next()
    }
}

// 요청 ID가 포함된 로거 반환
func LoggerFromContext(c *gin.Context) *slog.Logger {
    requestID, _ := c.Get(string(RequestIDKey))
    return slog.Default().With("request_id", requestID)
}
