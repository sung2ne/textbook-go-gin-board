
func (r *userRepository) SearchByUsername(ctx context.Context, query string, limit int) ([]*domain.User, error) {
    var users []*domain.User
    err := r.db.WithContext(ctx).
        Where("username LIKE ?", query+"%").
        Limit(limit).
        Find(&users).Error
    return users, err
}
