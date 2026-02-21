package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"goboardapi/internal/auth"
	"goboardapi/internal/config"
	"goboardapi/internal/database"
	"goboardapi/internal/handler"
	"goboardapi/internal/repository"
	"goboardapi/internal/router"
	"goboardapi/internal/service"

	_ "goboardapi/docs"

	"github.com/gin-gonic/gin"
)

// @title Go Board API
// @version 1.0
// @description 게시판 API 서버 (환경변수 기반 설정)

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT 인증 토큰. 형식: Bearer {access_token}
func main() {
	// 환경변수에서 설정 로드 (Docker 배포를 위한 env 기반 설정)
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/board")
	jwtSecret := getEnv("JWT_SECRET", "dev-secret")
	port := getEnvInt("PORT", 8080)
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	ginMode := getEnv("GIN_MODE", "debug")

	gin.SetMode(ginMode)

	// DATABASE_URL → DatabaseConfig 변환
	_ = dbURL // env var로 받아두되, 기존 config 방식 사용
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("설정 로드 실패: %v", err)
	}

	db, err := database.Init(&cfg.Database)

	// Redis 연결 (선택적)
	redisClient, err := database.NewRedis(redisAddr, "", 0)
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
		jwtSecret,
		15*time.Minute,
		24*time.Hour,
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

	addr := fmt.Sprintf(":%d", port)
	log.Printf("서버 시작: http://localhost%s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
