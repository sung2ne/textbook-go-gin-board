package service

import (
	"context"
	"errors"
	"time"

	"goboardapi/internal/auth"
	"goboardapi/internal/domain"
	"goboardapi/internal/dto"
	"goboardapi/internal/repository"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService interface {
	Signup(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	Logout(ctx context.Context, userID uint) error
}

type authService struct {
	userRepo        repository.UserRepository
	passwordService *auth.PasswordService
	tokenService    *auth.TokenService
}

func NewAuthService(
	userRepo repository.UserRepository,
	passwordService *auth.PasswordService,
	tokenService *auth.TokenService,
) AuthService {
	return &authService{
		userRepo:        userRepo,
		passwordService: passwordService,
		tokenService:    tokenService,
	}
}

func (s *authService) Signup(ctx context.Context, req *dto.SignupRequest) (*dto.SignupResponse, error) {
	if err := s.passwordService.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := s.passwordService.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashedPassword,
		Username: req.Username,
		Role:     domain.RoleUser,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return &dto.SignupResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.passwordService.Compare(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.tokenService.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Username,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.Update(ctx, user)

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.tokenService.GetAccessExpiry(),
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Role:     string(user.Role),
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	userID, err := s.tokenService.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	newAccessToken, err := s.tokenService.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Username,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.tokenService.GetAccessExpiry(),
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	return s.tokenService.RevokeRefreshToken(ctx, userID)
}
