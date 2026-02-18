package dto

import (
	"time"

	"goboardapi/internal/domain"
)

// ProfileResponse 프로필 응답
type ProfileResponse struct {
	ID          uint       `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	Stats       *UserStats `json:"stats,omitempty"`
}

// UserStats 사용자 통계
type UserStats struct {
	PostCount    int64 `json:"post_count"`
	CommentCount int64 `json:"comment_count"`
	LikeCount    int64 `json:"like_count"`
}

// UpdateProfileRequest 프로필 수정 요청
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
}

// ChangePasswordRequest 비밀번호 변경 요청
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// WithdrawRequest 회원 탈퇴 요청
type WithdrawRequest struct {
	Password string `json:"password" binding:"required"`
	Reason   string `json:"reason"`
}

// ToProfileResponse User → ProfileResponse 변환
func ToProfileResponse(user *domain.User, stats *UserStats) *ProfileResponse {
	return &ProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		Role:        string(user.Role),
		CreatedAt:   user.CreatedAt,
		LastLoginAt: user.LastLoginAt,
		Stats:       stats,
	}
}
