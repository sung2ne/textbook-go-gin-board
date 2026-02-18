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
}

// NewRouter 생성자
func NewRouter(
	tokenService *auth.TokenService,
	tokenStore auth.TokenStore,
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler,
) *Router {
	return &Router{
		engine:         gin.Default(),
		tokenService:   tokenService,
		tokenStore:     tokenStore,
		authHandler:    authHandler,
		postHandler:    postHandler,
		commentHandler: commentHandler,
	}
}

// Setup 라우트 설정
func (r *Router) Setup() *gin.Engine {
	// API 버전 그룹
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

		// 게시글 라우트
		posts := v1.Group("/posts")
		{
			posts.GET("", r.postHandler.GetList)
			posts.GET("/cursor", r.postHandler.GetListByCursor)
			posts.POST("", r.postHandler.Create)
			posts.GET("/:id", r.postHandler.GetByID)
			posts.PUT("/:id", r.postHandler.Update)
			posts.DELETE("/:id", r.postHandler.Delete)

			// 댓글 라우트
			posts.GET("/:postId/comments", r.commentHandler.GetByPostID)
			posts.POST("/:postId/comments", r.commentHandler.Create)
			posts.PUT("/:postId/comments/:id", r.commentHandler.Update)
			posts.DELETE("/:postId/comments/:id", r.commentHandler.Delete)
		}
	}

	// 헬스 체크
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}
