
func (r *likeRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&domain.Like{}).
        Where("user_id = ?", userID).
        Count(&count).Error
    return count, err
}
