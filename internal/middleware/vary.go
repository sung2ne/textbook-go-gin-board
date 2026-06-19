package middleware

import "github.com/gin-gonic/gin"

// Vary - 캐시 구분 헤더
func Vary(headers ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        for _, h := range headers {
            c.Header("Vary", h)
        }
        c.Next()
    }
}
