package handler

import (
    "net/http"
    "runtime/trace"

    "github.com/gin-gonic/gin"
)

func ListPosts(c *gin.Context) {
    ctx := c.Request.Context()

    // 데이터베이스 조회 구간
    ctx, task := trace.NewTask(ctx, "ListPosts")
    defer task.End()

    trace.WithRegion(ctx, "database query", func() {
        // DB 조회
    })

    trace.WithRegion(ctx, "serialize response", func() {
        // JSON 직렬화
    })

    c.JSON(http.StatusOK, gin.H{"posts": []string{}})
}
