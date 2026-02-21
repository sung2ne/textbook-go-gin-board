package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "goboardapi/internal/cache"
)

func CacheStats(c *gin.Context) {
    hits, misses, hitRate := cache.GetStats()

    c.JSON(http.StatusOK, gin.H{
        "hits":     hits,
        "misses":   misses,
        "hit_rate": hitRate,
    })
}
