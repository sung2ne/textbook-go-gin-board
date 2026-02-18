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

	var tokenStore auth.TokenStore
	if redisClient != nil {
		tokenStore = auth.NewRedisTokenStore(redisClient)
	}

	// 인증 서비스
	passwordService := auth.NewPasswordService()
	passwordHasher := auth.NewBcryptHasher(auth.DefaultCost)
	tokenService := auth.NewTokenService(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.AccessExpiry)*time.Minute,
		time.Duration(cfg.JWT.RefreshExpiry)*time.Hour,
		tokenStore,
	)

	// Repository
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	// Service
	authService := service.NewAuthService(userRepo, passwordService, tokenService)
	notificationSvc := service.NewNotificationService(notificationRepo, userRepo)
	postService := service.NewPostService(postRepo, likeRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, notificationSvc)
	likeService := service.NewLikeService(likeRepo, postRepo)
	adminService := service.NewAdminService(userRepo, postRepo, commentRepo)
	userService := service.NewUserService(userRepo, postRepo, commentRepo, likeRepo, passwordHasher, db)

	// Handler
	authHandler := handler.NewAuthHandler(authService, tokenService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	likeHandler := handler.NewLikeHandler(likeService)
	notificationHandler := handler.NewNotificationHandler(notificationSvc)
	adminHandler := handler.NewAdminHandler(adminService)
	userHandler := handler.NewUserHandler(userService)

	// 라우터 설정
	r := router.NewRouter(
		tokenService, tokenStore,
		authHandler, postHandler, commentHandler,
		likeHandler, notificationHandler, adminHandler,
		userHandler,
	)
	engine := r.Setup()

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("서버 시작: http://localhost%s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatal(err)
	}
}
