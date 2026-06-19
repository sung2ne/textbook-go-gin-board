
func (h *CommentHandler) UpdateComment(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 댓글 ID"})
        return
    }

    var req dto.UpdateCommentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    comment, err := h.commentService.Update(c.Request.Context(), uint(id), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToCommentResponse(comment)))
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 댓글 ID"})
        return
    }

    err = h.commentService.Delete(c.Request.Context(), uint(id))
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "댓글이 삭제되었습니다",
    })
}
