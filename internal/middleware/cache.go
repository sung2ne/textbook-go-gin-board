package middleware

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
)

// CacheControl - 캐시 헤더 설정
func CacheControl(maxAge time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(maxAge.Seconds())))
        c.Next()
    }
}

// NoCache - 캐시 비활성화
func NoCache() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
        c.Header("Pragma", "no-cache")
        c.Header("Expires", "0")
        c.Next()
    }
}
