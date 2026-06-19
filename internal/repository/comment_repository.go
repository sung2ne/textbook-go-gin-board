
func (r *commentRepository) CountByAuthorID(ctx context.Context, authorID uint) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&domain.Comment{}).
        Where("author_id = ?", authorID).
        Count(&count).Error
    return count, err
}

func (r *commentRepository) FindByAuthorIDWithPaging(ctx context.Context, authorID uint, offset, limit int) ([]*domain.Comment, int64, error) {
    var comments []*domain.Comment
    var total int64

    r.db.WithContext(ctx).
        Model(&domain.Comment{}).
        Where("author_id = ?", authorID).
        Count(&total)

    err := r.db.WithContext(ctx).
        Preload("Post").
        Where("author_id = ?", authorID).
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&comments).Error

    return comments, total, err
}
