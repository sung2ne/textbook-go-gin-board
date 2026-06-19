package handler

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func GetPost(c *gin.Context) {
    id := c.Param("id")
    post, err := postService.GetPost(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
        return
    }

    // Last-Modified 설정
    modTime := post.UpdatedAt.UTC()
    c.Header("Last-Modified", modTime.Format(http.TimeFormat))

    // If-Modified-Since 확인
    ifModifiedSince := c.GetHeader("If-Modified-Since")
    if ifModifiedSince != "" {
        t, err := time.Parse(http.TimeFormat, ifModifiedSince)
        if err == nil && !modTime.After(t) {
            c.Status(http.StatusNotModified)
            return
        }
    }

    c.JSON(http.StatusOK, post)
}
