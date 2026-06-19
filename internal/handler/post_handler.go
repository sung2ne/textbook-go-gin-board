
func (h *PostHandler) DeletePost(c *gin.Context) {
    claims := middleware.MustGetCurrentUser(c)
    id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

    post, err := h.postService.GetByID(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "게시글을 찾을 수 없습니다"})
        return
    }

    userRole := domain.Role(claims.Role)

    // 조건 1: 본인 게시글
    isOwner := post.AuthorID == claims.UserID

    // 조건 2: 관리 권한
    canManage := domain.HasPermission(userRole, domain.PermissionPostManage)

    if !isOwner && !canManage {
        c.JSON(http.StatusForbidden, gin.H{"error": "삭제 권한이 없습니다"})
        return
    }

    // 삭제 처리...
}
