func main() {
    cfg := config.Load()
    r := gin.Default()

    // 개발 환경에서만 Swagger 활성화
    if cfg.Environment != "production" {
        r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
        log.Printf("Swagger UI: http://%s/swagger/index.html", cfg.ServerAddress)
    }

    // ... 나머지 코드
}
