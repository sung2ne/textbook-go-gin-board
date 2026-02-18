package repository

import (
	"context"
	"errors"

	"goboardapi/internal/domain"

	"gorm.io/gorm"
)

var (
	ErrAlreadyLiked = errors.New("already liked")
	ErrNotLiked     = errors.New("not liked")
)

type LikeRepository interface {
	Like(ctx context.Context, userID, postID uint) error
	Unlike(ctx context.Context, userID, postID uint) error
	IsLiked(ctx context.Context, userID, postID uint) (bool, error)
	GetLikeCount(ctx context.Context, postID uint) (int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Like(ctx context.Context, userID, postID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 이미 좋아요 했는지 확인
		var count int64
		tx.Model(&domain.Like{}).
			Where("user_id = ? AND post_id = ?", userID, postID).
			Count(&count)

		if count > 0 {
			return ErrAlreadyLiked
		}

		// 좋아요 생성
		like := &domain.Like{
			UserID: userID,
			PostID: postID,
		}
		if err := tx.Create(like).Error; err != nil {
			return err
		}

		// 게시글의 좋아요 수 증가
		return tx.Model(&domain.Post{}).
			Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1")).
			Error
	})
}

func (r *likeRepository) Unlike(ctx context.Context, userID, postID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ? AND post_id = ?", userID, postID).
			Delete(&domain.Like{})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotLiked
		}

		// 게시글의 좋아요 수 감소
		return tx.Model(&domain.Post{}).
			Where("id = ?", postID).
			UpdateColumn("like_count", gorm.Expr("like_count - 1")).
			Error
	})
}

func (r *likeRepository) IsLiked(ctx context.Context, userID, postID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

func (r *likeRepository) GetLikeCount(ctx context.Context, postID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("post_id = ?", postID).
		Count(&count).Error
	return count, err
}

func (r *likeRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
