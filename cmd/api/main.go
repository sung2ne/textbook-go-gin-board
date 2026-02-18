package main

import (
	"fmt"
	"log"

	"goboardapi/internal/config"
	"goboardapi/internal/database"
	"goboardapi/internal/handler"
	"goboardapi/internal/repository"
	"goboardapi/internal/router"
	"goboardapi/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 설정 로드
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Gin 모드 설정
	gin.SetMode(cfg.Server.Mode)

	// 데이터베이스 연결
	db, err := database.Init(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	// 의존성 주입
	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo, cfg)
	postHandler := handler.NewPostHandler(postService)

	// 라우터 설정
	r := router.NewRouter(postHandler)
	engine := r.Setup()

	// 서버 시작
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("서버 시작: http://localhost%s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatal(err)
	}
}
