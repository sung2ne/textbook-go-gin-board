
func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&domain.User{}).
        Where("username = ?", username).
        Count(&count).Error
    return count > 0, err
}
