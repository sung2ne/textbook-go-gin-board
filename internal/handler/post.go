func (h *PostHandler) Create(c *gin.Context) {
    log := middleware.LoggerFromContext(c)

    log.Info("게시글 생성 시작")

    // ... 게시글 생성 로직

    log.Info("게시글 생성 완료", "post_id", post.ID)
}
