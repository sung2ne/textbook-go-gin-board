package router

import (
	"goboardapi/internal/auth"
	"goboardapi/internal/handler"
	"goboardapi/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router 라우터
type Router struct {
	engine         *gin.Engine
	tokenService   *auth.TokenService
	tokenStore     auth.TokenStore
	authHandler    *handler.AuthHandler
	postHandler    *handler.PostHandler
	commentHandler *handler.CommentHandler
	likeHandler    *handler.LikeHandler
}

// NewRouter 생성자
func NewRouter(
	tokenService *auth.TokenService,
	tokenStore auth.TokenStore,
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler,
	likeHandler *handler.LikeHandler,
) *Router {
	return &Router{
		engine:         gin.Default(),
		tokenService:   tokenService,
		tokenStore:     tokenStore,
		authHandler:    authHandler,
		postHandler:    postHandler,
		commentHandler: commentHandler,
		likeHandler:    likeHandler,
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
		}
	}

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}
