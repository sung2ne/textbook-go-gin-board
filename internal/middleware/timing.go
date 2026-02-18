package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Timing() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		c.Header("X-Response-Time", fmt.Sprintf("%dms", duration.Milliseconds()))
		c.Set("latency", duration)
	}
}
