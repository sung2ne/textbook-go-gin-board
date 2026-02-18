package main

import (
	"fmt"
	"log"
	"time"

	"goboardapi/internal/auth"
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

	// Redis 연결
	redisClient, err := database.NewRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Printf("Redis 연결 실패 (토큰 저장소 비활성화): %v", err)
	}

	// 토큰 저장소
	var tokenStore auth.TokenStore
	if redisClient != nil {
		tokenStore = auth.NewRedisTokenStore(redisClient)
	}

	// 인증 서비스
	passwordService := auth.NewPasswordService()
	tokenService := auth.NewTokenService(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.AccessExpiry)*time.Minute,
		time.Duration(cfg.JWT.RefreshExpiry)*time.Hour,
		tokenStore,
	)

	// 의존성 주입
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	authService := service.NewAuthService(userRepo, passwordService, tokenService)
	postService := service.NewPostService(postRepo, cfg)
	commentService := service.NewCommentService(commentRepo, postRepo)

	authHandler := handler.NewAuthHandler(authService, tokenService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	// 라우터 설정
	r := router.NewRouter(tokenService, tokenStore, authHandler, postHandler, commentHandler)
	engine := r.Setup()

	// 서버 시작
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("서버 시작: http://localhost%s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatal(err)
	}
}
