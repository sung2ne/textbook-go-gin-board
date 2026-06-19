
// 좋아요 (인증 필요)
protected.POST("/posts/:id/like", likeHandler.LikePost)
protected.DELETE("/posts/:id/like", likeHandler.UnlikePost)
