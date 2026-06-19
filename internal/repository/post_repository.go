
func (r *postRepository) CountByAuthorID(ctx context.Context, authorID uint) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&domain.Post{}).
        Where("author_id = ?", authorID).
        Count(&count).Error
    return count, err
}
