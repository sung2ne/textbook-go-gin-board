package service

import (
	"context"

	"goboardapi/internal/domain"
	"goboardapi/internal/dto"
	"goboardapi/internal/middleware"
	"goboardapi/internal/repository"
)

type AdminService interface {
	GetStats(ctx context.Context) (*dto.AdminStats, error)
	ListUsers(ctx context.Context, page, size int) ([]*dto.UserListResponse, int64, error)
	ChangeRole(ctx context.Context, userID uint, req *dto.ChangeRoleRequest) error
	DeleteUser(ctx context.Context, userID uint) error
	ForceDeletePost(ctx context.Context, postID uint) error
}

type adminService struct {
	userRepo    repository.UserRepository
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
}

func NewAdminService(
	userRepo repository.UserRepository,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
) AdminService {
	return &adminService{
		userRepo:    userRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

func (s *adminService) GetStats(ctx context.Context) (*dto.AdminStats, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}
	if claims.Role != string(domain.RoleAdmin) {
		return nil, ErrForbidden
	}

	userCount, _ := s.userRepo.Count(ctx)
	postCount, _ := s.postRepo.CountByAuthorID(ctx, 0) // 0 for all
	commentCount, _ := s.commentRepo.CountByAuthorID(ctx, 0)

	return &dto.AdminStats{
		TotalUsers:    userCount,
		TotalPosts:    postCount,
		TotalComments: commentCount,
	}, nil
}

func (s *adminService) ListUsers(ctx context.Context, page, size int) ([]*dto.UserListResponse, int64, error) {
	offset := (page - 1) * size
	users, total, err := s.userRepo.FindAll(ctx, offset, size)
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.UserListResponse
	for _, user := range users {
		responses = append(responses, &dto.UserListResponse{
			ID:          user.ID,
			Email:       user.Email,
			Username:    user.Username,
			Role:        string(user.Role),
			LastLoginAt: user.LastLoginAt,
			CreatedAt:   user.CreatedAt,
		})
	}

	return responses, total, nil
}

func (s *adminService) ChangeRole(ctx context.Context, userID uint, req *dto.ChangeRoleRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Role = domain.Role(req.Role)
	return s.userRepo.Update(ctx, user)
}

func (s *adminService) DeleteUser(ctx context.Context, userID uint) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *adminService) ForceDeletePost(ctx context.Context, postID uint) error {
	return s.postRepo.Delete(ctx, postID)
}
