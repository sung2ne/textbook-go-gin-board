
// GetList 게시글 목록 조회
// GET /api/v1/posts?page=1&size=10&q=검색어&type=title&sort=views,desc
func (h *PostHandler) GetList(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

    search := &dto.SearchParams{
        Query:      c.Query("q"),
        SearchType: c.Query("type"),
    }

    sort := &dto.SortParams{
        Sort: c.Query("sort"),
    }

    posts, meta, err := h.postService.GetList(page, size, search, sort)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse("SERVER_ERROR", "목록 조회에 실패했습니다"))
        return
    }

    c.JSON(http.StatusOK, dto.SuccessWithMeta(posts, meta))
}
