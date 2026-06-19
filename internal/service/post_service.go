
// GetList 게시글 목록 조회 (검색, 정렬 포함)
func (s *PostService) GetList(page, size int, search *dto.SearchParams, sort *dto.SortParams) ([]dto.PostListResponse, *dto.Meta, error) {
    pagination := dto.NewPagination(
        page,
        size,
        s.cfg.Pagination.DefaultSize,
        s.cfg.Pagination.MaxSize,
    )

    posts, total, err := s.postRepo.FindAll(pagination, search, sort)
    if err != nil {
        return nil, nil, err
    }

    list := make([]dto.PostListResponse, len(posts))
    for i, post := range posts {
        list[i] = dto.PostListResponse{
            ID:        post.ID,
            Title:     post.Title,
            Author:    post.Author,
            Views:     post.Views,
            CreatedAt: post.CreatedAt,
        }
    }

    totalPages := int(total) / pagination.Size
    if int(total)%pagination.Size > 0 {
        totalPages++
    }

    meta := &dto.Meta{
        Page:       pagination.Page,
        Size:       pagination.Size,
        Total:      total,
        TotalPages: totalPages,
    }

    return list, meta, nil
}
