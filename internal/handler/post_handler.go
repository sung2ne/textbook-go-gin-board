
func (h *PostHandler) UpdatePost(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID 형식"})
        return
    }

    var req dto.UpdatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    post, err := h.postService.Update(c.Request.Context(), uint(id), &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, dto.SuccessResponse(dto.ToPostResponse(post)))
}

func (h *PostHandler) DeletePost(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 ID 형식"})
        return
    }

    err = h.postService.Delete(c.Request.Context(), uint(id))
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "게시글이 삭제되었습니다",
    })
}
