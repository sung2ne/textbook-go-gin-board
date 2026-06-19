
func (r *userRepository) FindInactiveBatch(ctx context.Context, afterID uint, limit int) ([]*domain.User, error) {
    var users []*domain.User

    thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

    err := r.db.WithContext(ctx).
        Where("id > ?", afterID).
        Where("last_login_at < ?", thirtyDaysAgo).
        Order("id ASC").
        Limit(limit).
        Find(&users).Error

    return users, err
}
