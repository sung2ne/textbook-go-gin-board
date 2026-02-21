package repository

import (
	"context"
	"errors"

	"goboardapi/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	FindByID(ctx context.Context, id uint) (*domain.Post, error)
	FindAll(ctx context.Context, offset, limit int) ([]*domain.Post, int64, error)
	Update(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, id uint) error
	FindByAuthorID(ctx context.Context, authorID uint, offset, limit int) ([]*domain.Post, int64, error)
	CountByAuthorID(ctx context.Context, authorID uint) (int64, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) FindByID(ctx context.Context, id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.WithContext(ctx).
		Preload("Author").
		First(&post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindAll(ctx context.Context, offset, limit int) ([]*domain.Post, int64, error) {
	var posts []*domain.Post
	var total int64

	r.db.WithContext(ctx).Model(&domain.Post{}).Count(&total)

	err := r.db.WithContext(ctx).
		Preload("Author").
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *postRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Post{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPostNotFound
	}
	return nil
}

func (r *postRepository) FindByAuthorID(ctx context.Context, authorID uint, offset, limit int) ([]*domain.Post, int64, error) {
	var posts []*domain.Post
	var total int64

	r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Where("author_id = ?", authorID).
		Count(&total)

	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("author_id = ?", authorID).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) CountByAuthorID(ctx context.Context, authorID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Where("author_id = ?", authorID).
		Count(&count).Error
	return count, err
}
