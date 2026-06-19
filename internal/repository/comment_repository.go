
func (r *commentRepository) HasReplies(ctx context.Context, parentID uint) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&domain.Comment{}).
        Where("parent_id = ?", parentID).
        Count(&count).Error
    return count > 0, err
}
