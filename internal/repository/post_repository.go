package repository

type PostRepository struct {
    db *gorm.DB
}

func (r *PostRepository) Create(post *domain.Post) error {
    return r.db.Create(post).Error
}
