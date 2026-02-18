package service

import (
	"context"

	"goboardapi/internal/auth"
	"goboardapi/internal/domain"
	"goboardapi/internal/dto"
	"goboardapi/internal/middleware"
	"goboardapi/internal/repository"
)

type PostService interface {
	Create(ctx context.Context, req *dto.CreatePostRequest) (*domain.Post, error)
	GetByID(ctx context.Context, id uint) (*dto.PostResponse, error)
	GetAll(ctx context.Context, page, size int) ([]*domain.Post, int64, error)
	Update(ctx context.Context, id uint, req *dto.UpdatePostRequest) (*domain.Post, error)
	Delete(ctx context.Context, id uint) error
}

type postService struct {
	postRepo repository.PostRepository
	likeRepo repository.LikeRepository
}

func NewPostService(postRepo repository.PostRepository, likeRepo repository.LikeRepository) PostService {
	return &postService{
		postRepo: postRepo,
		likeRepo: likeRepo,
	}
}

func (s *postService) Create(ctx context.Context, req *dto.CreatePostRequest) (*domain.Post, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	post := &domain.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: claims.UserID,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	return s.postRepo.FindByID(ctx, post.ID)
}

func (s *postService) GetByID(ctx context.Context, id uint) (*dto.PostResponse, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := dto.ToPostResponse(post)

	// 로그인 상태면 좋아요 여부 확인
	if claims, ok := middleware.GetUserFromContext(ctx); ok {
		isLiked, _ := s.likeRepo.IsLiked(ctx, claims.UserID, id)
		resp.IsLiked = isLiked
	}

	return resp, nil
}

func (s *postService) GetAll(ctx context.Context, page, size int) ([]*domain.Post, int64, error) {
	offset := (page - 1) * size
	return s.postRepo.FindAll(ctx, offset, size)
}

func (s *postService) Update(ctx context.Context, id uint, req *dto.UpdatePostRequest) (*domain.Post, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !s.canModify(post.AuthorID, claims) {
		return nil, ErrForbidden
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := s.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) Delete(ctx context.Context, id uint) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if !s.canModify(post.AuthorID, claims) {
		return ErrForbidden
	}

	return s.postRepo.Delete(ctx, id)
}

// canModify 수정/삭제 권한 확인
func (s *postService) canModify(authorID uint, claims *auth.CustomClaims) bool {
	if authorID == claims.UserID {
		return true
	}
	if claims.Role == string(domain.RoleAdmin) {
		return true
	}
	return false
}
