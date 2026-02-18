package repository

import (
	"goboardapi/internal/domain"

	"gorm.io/gorm"
)

// PostRepository 게시글 저장소 인터페이스
type PostRepository interface {
	Create(post *domain.Post) error
	FindByID(id uint) (*domain.Post, error)
	FindAll(offset, limit int) ([]domain.Post, int64, error)
	Update(post *domain.Post) error
	Delete(id uint) error
	IncrementViews(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindAll(offset, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	if err := r.db.Model(&domain.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Post{}, id).Error
}

func (r *postRepository) IncrementViews(id uint) error {
	return r.db.Model(&domain.Post{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1")).
		Error
}
