package service

import (
	"context"

	"goboardapi/internal/auth"
	"goboardapi/internal/domain"
	"goboardapi/internal/dto"
	"goboardapi/internal/middleware"
	"goboardapi/internal/repository"

	"gorm.io/gorm"
)

type UserService interface {
	GetProfile(ctx context.Context) (*dto.ProfileResponse, error)
	UpdateProfile(ctx context.Context, req *dto.UpdateProfileRequest) (*dto.ProfileResponse, error)
	ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) error
	Withdraw(ctx context.Context, req *dto.WithdrawRequest) error
	GetMyPosts(ctx context.Context, page, size int) ([]*domain.Post, int64, error)
	GetMyComments(ctx context.Context, page, size int) ([]*domain.Comment, int64, error)
	SearchByUsername(ctx context.Context, query string, limit int) ([]*domain.User, error)
}

type userService struct {
	userRepo       repository.UserRepository
	postRepo       repository.PostRepository
	commentRepo    repository.CommentRepository
	likeRepo       repository.LikeRepository
	passwordHasher auth.PasswordHasher
	db             *gorm.DB
}

func NewUserService(
	userRepo repository.UserRepository,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	likeRepo repository.LikeRepository,
	passwordHasher auth.PasswordHasher,
	db *gorm.DB,
) UserService {
	return &userService{
		userRepo:       userRepo,
		postRepo:       postRepo,
		commentRepo:    commentRepo,
		likeRepo:       likeRepo,
		passwordHasher: passwordHasher,
		db:             db,
	}
}

func (s *userService) GetProfile(ctx context.Context) (*dto.ProfileResponse, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	stats, err := s.getUserStats(ctx, claims.UserID)
	if err != nil {
		stats = nil
	}

	return dto.ToProfileResponse(user, stats), nil
}

func (s *userService) UpdateProfile(ctx context.Context, req *dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if req.Username != user.Username {
		exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrUsernameExists
		}
	}

	user.Username = req.Username

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	stats, _ := s.getUserStats(ctx, claims.UserID)

	return dto.ToProfileResponse(user, stats), nil
}

func (s *userService) ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return err
	}

	if !s.passwordHasher.Compare(user.Password, req.CurrentPassword) {
		return ErrWrongPassword
	}

	if s.passwordHasher.Compare(user.Password, req.NewPassword) {
		return ErrSamePassword
	}

	hashedPassword, err := s.passwordHasher.Hash(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(ctx, user)
}

func (s *userService) Withdraw(ctx context.Context, req *dto.WithdrawRequest) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return err
	}

	if user.Role == domain.RoleAdmin {
		return ErrCannotWithdrawAdmin
	}

	if !s.passwordHasher.Compare(user.Password, req.Password) {
		return ErrWrongPassword
	}

	return s.withdraw(ctx, user, req.Reason)
}

func (s *userService) withdraw(ctx context.Context, user *domain.User, reason string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 좋아요 삭제
		if err := tx.Where("user_id = ?", user.ID).Delete(&domain.Like{}).Error; err != nil {
			return err
		}

		// 알림 삭제
		if err := tx.Where("user_id = ?", user.ID).Delete(&domain.Notification{}).Error; err != nil {
			return err
		}

		// 게시글 작성자 익명화
		if err := tx.Model(&domain.Post{}).
			Where("author_id = ?", user.ID).
			Update("author_id", 0).Error; err != nil {
			return err
		}

		// 댓글 작성자 익명화
		if err := tx.Model(&domain.Comment{}).
			Where("author_id = ?", user.ID).
			Update("author_id", 0).Error; err != nil {
			return err
		}

		// 탈퇴 사유 저장
		if reason != "" {
			withdrawLog := &domain.WithdrawLog{
				UserID: user.ID,
				Email:  user.Email,
				Reason: reason,
			}
			if err := tx.Create(withdrawLog).Error; err != nil {
				return err
			}
		}

		// 사용자 소프트 삭제
		return tx.Delete(user).Error
	})
}

func (s *userService) GetMyPosts(ctx context.Context, page, size int) ([]*domain.Post, int64, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, 0, ErrUnauthorized
	}

	offset := (page - 1) * size
	return s.postRepo.FindByAuthorID(ctx, claims.UserID, offset, size)
}

func (s *userService) GetMyComments(ctx context.Context, page, size int) ([]*domain.Comment, int64, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, 0, ErrUnauthorized
	}

	offset := (page - 1) * size
	return s.commentRepo.FindByAuthorID(ctx, claims.UserID, offset, size)
}

func (s *userService) SearchByUsername(ctx context.Context, query string, limit int) ([]*domain.User, error) {
	return s.userRepo.SearchByUsername(ctx, query, limit)
}

func (s *userService) getUserStats(ctx context.Context, userID uint) (*dto.UserStats, error) {
	postCount, err := s.postRepo.CountByAuthorID(ctx, userID)
	if err != nil {
		return nil, err
	}

	commentCount, err := s.commentRepo.CountByAuthorID(ctx, userID)
	if err != nil {
		return nil, err
	}

	likeCount, err := s.likeRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserStats{
		PostCount:    postCount,
		CommentCount: commentCount,
		LikeCount:    likeCount,
	}, nil
}
