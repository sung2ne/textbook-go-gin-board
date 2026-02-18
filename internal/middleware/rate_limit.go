package middleware

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "요청이 너무 많습니다. 잠시 후 다시 시도해주세요.",
				},
			})
			return
		}
		c.Next()
	}
}

type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     sync.RWMutex
	r      rate.Limit
	b      int
	expiry time.Duration
}

func NewIPRateLimiter(r rate.Limit, b int, expiry time.Duration) *IPRateLimiter {
	irl := &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		r:      r,
		b:      b,
		expiry: expiry,
	}
	go irl.cleanupLoop()
	return irl
}

func (irl *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	irl.mu.Lock()
	defer irl.mu.Unlock()

	limiter, exists := irl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(irl.r, irl.b)
		irl.ips[ip] = limiter
	}
	return limiter
}

func (irl *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(irl.expiry)
	for range ticker.C {
		irl.mu.Lock()
		irl.ips = make(map[string]*rate.Limiter)
		irl.mu.Unlock()
	}
}

func (irl *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := irl.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "요청이 너무 많습니다. 잠시 후 다시 시도해주세요.",
				},
			})
			return
		}
		c.Next()
	}
}

type RateLimitConfig struct {
	Rate   rate.Limit
	Burst  int
	Window time.Duration
}

func RateLimiterWithHeaders(cfg RateLimitConfig) gin.HandlerFunc {
	limiter := rate.NewLimiter(cfg.Rate, cfg.Burst)

	return func(c *gin.Context) {
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Burst))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%.0f", limiter.Tokens()))

		if !limiter.Allow() {
			reservation := limiter.Reserve()
			delay := reservation.Delay()
			reservation.Cancel()

			c.Header("Retry-After", fmt.Sprintf("%d", int(delay.Seconds())+1))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":        "RATE_LIMITED",
					"message":     "요청이 너무 많습니다.",
					"retry_after": int(delay.Seconds()) + 1,
				},
			})
			return
		}
		c.Next()
	}
}

type RedisRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRedisRateLimiter(client *redis.Client, limit int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

func (rl *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := rl.client.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rl.client.Expire(ctx, key, rl.window)
		}

		if count > int64(rl.limit) {
			ttl, _ := rl.client.TTL(ctx, key).Result()
			c.Header("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())+1))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "요청이 너무 많습니다.",
				},
			})
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.limit-int(count)))
		c.Next()
	}
}
