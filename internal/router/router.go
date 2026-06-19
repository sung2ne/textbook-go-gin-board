
type Router struct {
    engine         *gin.Engine
    postHandler    *handler.PostHandler
    commentHandler *handler.CommentHandler
}

func NewRouter(postHandler *handler.PostHandler, commentHandler *handler.CommentHandler) *Router {
    return &Router{
        engine:         gin.Default(),
        postHandler:    postHandler,
        commentHandler: commentHandler,
    }
}

func (r *Router) Setup() *gin.Engine {
    v1 := r.engine.Group("/api/v1")
    {
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

    r.engine.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    return r.engine
}
