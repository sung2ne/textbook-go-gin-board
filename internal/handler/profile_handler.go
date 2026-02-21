package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "goboardapi/internal/service"
)

func SlowEndpoint(c *gin.Context) {
    input := c.DefaultQuery("input", "hello")
    result := service.SlowHash(input)
    c.JSON(http.StatusOK, gin.H{"hash": result})
}

func ConcatEndpoint(c *gin.Context) {
    items := make([]string, 10000)
    for i := range items {
        items[i] = "item"
    }

    result := service.InefficientConcat(items)
    c.JSON(http.StatusOK, gin.H{"length": len(result)})
}
