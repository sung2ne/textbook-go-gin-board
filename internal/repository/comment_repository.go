package repository

import (
	"goboardapi/internal/domain"

	"gorm.io/gorm"
)

// CommentRepository 댓글 저장소 인터페이스
type CommentRepository interface {
	Create(comment *domain.Comment) error
	FindByID(id uint) (*domain.Comment, error)
	FindByPostID(postID uint) ([]domain.Comment, error)
	FindByPostIDWithReplies(postID uint) ([]domain.Comment, error)
	Update(comment *domain.Comment) error
	Delete(id uint) error
	HasReplies(id uint) (bool, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByID(id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByPostID(postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) FindByPostIDWithReplies(postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment

	// 최상위 댓글 조회 (대댓글은 Preload로 함께 가져옴)
	err := r.db.
		Where("post_id = ? AND parent_id IS NULL", postID).
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Order("created_at ASC").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *commentRepository) Update(comment *domain.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Comment{}, id).Error
}

func (r *commentRepository) HasReplies(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Comment{}).Where("parent_id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
