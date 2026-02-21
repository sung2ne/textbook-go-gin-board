package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HealthCheck(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        sqlDB, _ := db.DB()
        ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
        defer cancel()

        if err := sqlDB.PingContext(ctx); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"status": "healthy"})
    }
}
