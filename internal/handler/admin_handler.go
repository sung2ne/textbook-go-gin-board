package handler

// InvalidateCache - 캐시 무효화 API
func InvalidateCache(c *gin.Context) {
    path := c.Query("path")

    // CDN API 호출 (CloudFront, Cloudflare 등)
    // cdn.Invalidate(path)

    // Surrogate-Key로 선택적 무효화
    c.Header("Surrogate-Key", "posts")

    c.JSON(http.StatusOK, gin.H{"invalidated": path})
}
