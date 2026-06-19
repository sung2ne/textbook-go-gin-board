
// PostRepository 인터페이스 수정
type PostRepository interface {
    Create(post *domain.Post) error
    FindByID(id uint) (*domain.Post, error)
    FindAll(pagination *dto.Pagination, search *dto.SearchParams, sort *dto.SortParams) ([]domain.Post, int64, error)
    Update(post *domain.Post) error
    Delete(id uint) error
    IncrementViews(id uint) error
    FindAllByCursor(cursor *dto.Cursor, limit int) ([]domain.Post, error)
}

// FindAll 구현 수정
func (r *postRepository) FindAll(pagination *dto.Pagination, search *dto.SearchParams, sort *dto.SortParams) ([]domain.Post, int64, error) {
    var posts []domain.Post
    var total int64

    query := r.db.Model(&domain.Post{})

    // 검색 조건 적용
    if search != nil && search.Query != "" {
        searchQuery := "%" + search.Query + "%"
        switch search.GetSearchType() {
        case dto.SearchTypeTitle:
            query = query.Where("title ILIKE ?", searchQuery)
        case dto.SearchTypeContent:
            query = query.Where("content ILIKE ?", searchQuery)
        default:
            query = query.Where("title ILIKE ? OR content ILIKE ?", searchQuery, searchQuery)
        }
    }

    // 전체 개수 조회
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // 정렬 조건 적용
    orderStr := "created_at DESC"
    if sort != nil {
        orderStr = sort.ToOrderString()
    }

    // 페이징 및 정렬 적용하여 조회
    err := query.
        Order(orderStr).
        Offset(pagination.Offset()).
        Limit(pagination.Size).
        Find(&posts).Error

    if err != nil {
        return nil, 0, err
    }

    return posts, total, nil
}
