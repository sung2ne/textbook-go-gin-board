package handler

func SetupRoutes(r *gin.Engine) {
    // 버전이 포함된 정적 파일 - 1년 캐시
    r.Static("/static", "./static")
    r.Use(func(c *gin.Context) {
        if strings.HasPrefix(c.Request.URL.Path, "/static") {
            c.Header("Cache-Control", "public, max-age=31536000, immutable")
        }
        c.Next()
    })
}
