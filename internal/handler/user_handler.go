
// SearchUsers 사용자 검색 (멘션 자동완성용)
func (h *UserHandler) SearchUsers(c *gin.Context) {
    query := c.Query("q")
    if len(query) < 2 {
        c.JSON(http.StatusOK, dto.SuccessResponse([]interface{}{}))
        return
    }

    users, err := h.userService.SearchByUsername(c.Request.Context(), query, 10)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 오류"})
        return
    }

    var results []gin.H
    for _, user := range users {
        results = append(results, gin.H{
            "id":       user.ID,
            "username": user.Username,
        })
    }

    c.JSON(http.StatusOK, dto.SuccessResponse(results))
}
