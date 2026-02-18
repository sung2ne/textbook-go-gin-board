package service

import (
	"context"

	"goboardapi/internal/middleware"
	"goboardapi/internal/repository"
)

type LikeService interface {
	Like(ctx context.Context, postID uint) error
	Unlike(ctx context.Context, postID uint) error
	IsLiked(ctx context.Context, postID uint) (bool, error)
}

type likeService struct {
	likeRepo repository.LikeRepository
	postRepo repository.PostRepository
}

func NewLikeService(likeRepo repository.LikeRepository, postRepo repository.PostRepository) LikeService {
	return &likeService{
		likeRepo: likeRepo,
		postRepo: postRepo,
	}
}

func (s *likeService) Like(ctx context.Context, postID uint) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	_, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return err
	}

	return s.likeRepo.Like(ctx, claims.UserID, postID)
}

func (s *likeService) Unlike(ctx context.Context, postID uint) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	return s.likeRepo.Unlike(ctx, claims.UserID, postID)
}

func (s *likeService) IsLiked(ctx context.Context, postID uint) (bool, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return false, nil
	}

	return s.likeRepo.IsLiked(ctx, claims.UserID, postID)
}
