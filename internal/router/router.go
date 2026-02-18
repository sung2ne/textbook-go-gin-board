package router

import (
	"goboardapi/internal/auth"
	"goboardapi/internal/handler"
	"goboardapi/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router 라우터
type Router struct {
	engine              *gin.Engine
	tokenService        *auth.TokenService
	tokenStore          auth.TokenStore
	authHandler         *handler.AuthHandler
	postHandler         *handler.PostHandler
	commentHandler      *handler.CommentHandler
	likeHandler         *handler.LikeHandler
	notificationHandler *handler.NotificationHandler
	adminHandler        *handler.AdminHandler
	userHandler         *handler.UserHandler
}

// NewRouter 생성자
func NewRouter(
	tokenService *auth.TokenService,
	tokenStore auth.TokenStore,
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler,
	likeHandler *handler.LikeHandler,
	notificationHandler *handler.NotificationHandler,
	adminHandler *handler.AdminHandler,
	userHandler *handler.UserHandler,
) *Router {
	return &Router{
		engine:              gin.Default(),
		tokenService:        tokenService,
		tokenStore:          tokenStore,
		authHandler:         authHandler,
		postHandler:         postHandler,
		commentHandler:      commentHandler,
		likeHandler:         likeHandler,
		notificationHandler: notificationHandler,
		adminHandler:        adminHandler,
		userHandler:         userHandler,
	}
}

// Setup 라우트 설정
func (r *Router) Setup() *gin.Engine {
	v1 := r.engine.Group("/api/v1")
	{
		// 인증 라우트
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/signup", r.authHandler.Signup)
			authGroup.POST("/login", r.authHandler.Login)
			authGroup.POST("/refresh", r.authHandler.RefreshToken)
			authGroup.POST("/logout",
				middleware.AuthMiddleware(r.tokenService, r.tokenStore),
				r.authHandler.Logout,
			)
		}

		// 공개 API
		posts := v1.Group("/posts")
		{
			posts.GET("", r.postHandler.GetList)
			posts.GET("/:id", r.postHandler.GetByID)
			posts.GET("/:postId/comments", r.commentHandler.GetByPostID)
		}

		// 인증 필요 API
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(r.tokenService, r.tokenStore))
		{
			// 게시글 작성/수정/삭제
			protected.POST("/posts", r.postHandler.Create)
			protected.PUT("/posts/:id", r.postHandler.Update)
			protected.DELETE("/posts/:id", r.postHandler.Delete)

			// 좋아요
			protected.POST("/posts/:id/like", r.likeHandler.LikePost)
			protected.DELETE("/posts/:id/like", r.likeHandler.UnlikePost)

			// 댓글 작성/수정/삭제
			protected.POST("/posts/:postId/comments", r.commentHandler.Create)
			protected.PUT("/comments/:id", r.commentHandler.Update)
			protected.DELETE("/comments/:id", r.commentHandler.Delete)

			// 알림
			protected.GET("/notifications", r.notificationHandler.GetNotifications)
			protected.GET("/notifications/unread-count", r.notificationHandler.GetUnreadCount)
			protected.PUT("/notifications/:id/read", r.notificationHandler.MarkAsRead)
			protected.PUT("/notifications/read-all", r.notificationHandler.MarkAllAsRead)

			// 사용자 검색 (멘션 자동완성)
			protected.GET("/users/search", r.userHandler.SearchUsers)

			// 마이페이지
			me := protected.Group("/me")
			{
				me.GET("/profile", r.userHandler.GetProfile)
				me.PUT("/profile", r.userHandler.UpdateProfile)
				me.PUT("/password", r.userHandler.ChangePassword)
				me.DELETE("", r.userHandler.Withdraw)
				me.GET("/posts", r.userHandler.GetMyPosts)
				me.GET("/comments", r.userHandler.GetMyComments)
			}

			// 관리자 API
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/stats", r.adminHandler.GetStats)
				admin.GET("/users", r.adminHandler.ListUsers)
				admin.PUT("/users/:id/role", r.adminHandler.ChangeRole)
				admin.DELETE("/users/:id", r.adminHandler.DeleteUser)
				admin.DELETE("/posts/:id", r.adminHandler.ForceDeletePost)
			}
		}
	}

	// Swagger UI
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.DeepLinking(true),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.DocExpansion("list"),
	))

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}
