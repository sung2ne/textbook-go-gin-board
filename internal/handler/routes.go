package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"goboardapi/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// 정적 파일 - 1시간 캐시
	r.Static("/static", "./static")
	r.Use(middleware.CacheControl(1 * time.Hour))

	// API 그룹
	api := r.Group("/api")
	{
		// 공개 데이터 - 5분 캐시
		public := api.Group("/public")
		public.Use(middleware.CacheControl(5 * time.Minute))
		public.GET("/posts", ListPosts)
	}
}
