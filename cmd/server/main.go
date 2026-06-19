package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "go-board/docs"
    "go-board/internal/config"
)

// @title Go Board API
// @version 1.0
// @description 게시판 API 서버입니다.

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT 인증 토큰

// @tag.name auth
// @tag.description 인증 관련 API

// @tag.name posts
// @tag.description 게시글 CRUD API

// @tag.name comments
// @tag.description 댓글 CRUD API

func main() {
    cfg := config.Load()

    // Swagger 동적 설정
    if host := os.Getenv("SWAGGER_HOST"); host != "" {
        docs.SwaggerInfo.Host = host
    } else {
        docs.SwaggerInfo.Host = cfg.ServerAddress
    }

    r := gin.Default()

    // 프로덕션이 아닐 때만 Swagger 활성화
    if cfg.Environment != "production" {
        r.GET("/swagger/*any", ginSwagger.WrapHandler(
            swaggerFiles.Handler,
            ginSwagger.DeepLinking(true),
            ginSwagger.DefaultModelsExpandDepth(-1),
            ginSwagger.DocExpansion("list"),
            ginSwagger.PersistAuthorization(cfg.Environment == "development"),
        ))
        log.Printf("Swagger UI: http://%s/swagger/index.html", cfg.ServerAddress)
    }

    // 라우트 설정...
    setupRoutes(r, cfg)

    r.Run(cfg.ServerAddress)
}
