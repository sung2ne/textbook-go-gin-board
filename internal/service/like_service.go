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

    // 게시글 존재 확인
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
        return false, nil // 비로그인 상태면 좋아요 안함
    }

    return s.likeRepo.IsLiked(ctx, claims.UserID, postID)
}
