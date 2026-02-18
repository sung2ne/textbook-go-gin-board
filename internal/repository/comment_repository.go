package repository

import (
	"context"
	"errors"

	"goboardapi/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrCommentNotFound = errors.New("comment not found")
)

// CommentRepository 댓글 저장소 인터페이스
type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	FindByID(ctx context.Context, id uint) (*domain.Comment, error)
	FindByPostID(ctx context.Context, postID uint) ([]domain.Comment, error)
	FindByPostIDWithReplies(ctx context.Context, postID uint) ([]domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id uint) error
	HasReplies(ctx context.Context, parentID uint) (bool, error)
	FindByAuthorID(ctx context.Context, authorID uint, offset, limit int) ([]*domain.Comment, int64, error)
	CountByAuthorID(ctx context.Context, authorID uint) (int64, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) FindByID(ctx context.Context, id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		First(&comment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByPostID(ctx context.Context, postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) FindByPostIDWithReplies(ctx context.Context, postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment

	err := r.db.WithContext(ctx).
		Preload("Author").
		Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Author").Order("created_at ASC")
		}).
		Order("created_at ASC").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Comment{}, id).Error
}

func (r *commentRepository) HasReplies(ctx context.Context, parentID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("parent_id = ?", parentID).
		Count(&count).Error
	return count > 0, err
}

func (r *commentRepository) FindByAuthorID(ctx context.Context, authorID uint, offset, limit int) ([]*domain.Comment, int64, error) {
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

func (r *commentRepository) CountByAuthorID(ctx context.Context, authorID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("author_id = ?", authorID).
		Count(&count).Error
	return count, err
}
