package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

func (h *PostHandler) GetPost(c *gin.Context) {
    ctx := c.Request.Context() // HTTP 요청의 context

    id, _ := strconv.Atoi(c.Param("id"))

    post, err := h.service.FindByID(ctx, uint(id))
    if err != nil {
        if err == context.Canceled {
            // 클라이언트가 연결 끊음
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, post)
}
