package dto

import (
	"time"

	"goboardapi/internal/domain"
)

// CreatePostRequest 게시글 생성 요청
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required"`
}

// UpdatePostRequest 게시글 수정 요청
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,max=200"`
	Content string `json:"content" binding:"required"`
}

// AuthorInfo 작성자 정보
type AuthorInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// PostResponse 게시글 응답
type PostResponse struct {
	ID        uint        `json:"id"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Author    *AuthorInfo `json:"author"`
	Views     int         `json:"views"`
	LikeCount int         `json:"like_count"`
	IsLiked   bool        `json:"is_liked"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// PostListResponse 게시글 목록 응답
type PostListResponse struct {
	ID        uint        `json:"id"`
	Title     string      `json:"title"`
	Author    *AuthorInfo `json:"author"`
	Views     int         `json:"views"`
	CreatedAt time.Time   `json:"created_at"`
}

// ToPostResponse 도메인 → 응답 DTO 변환
func ToPostResponse(post *domain.Post) *PostResponse {
	resp := &PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Views:     post.Views,
		LikeCount: post.LikeCount,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	if post.AuthorID == 0 || post.Author == nil {
		resp.Author = &AuthorInfo{
			ID:       0,
			Username: "탈퇴한 사용자",
		}
	} else {
		resp.Author = &AuthorInfo{
			ID:       post.Author.ID,
			Username: post.Author.Username,
		}
	}

	return resp
}

// ToPostListResponse 목록용 변환
func ToPostListResponse(post *domain.Post) *PostListResponse {
	resp := &PostListResponse{
		ID:        post.ID,
		Title:     post.Title,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
	}

	if post.AuthorID == 0 || post.Author == nil {
		resp.Author = &AuthorInfo{
			ID:       0,
			Username: "탈퇴한 사용자",
		}
	} else {
		resp.Author = &AuthorInfo{
			ID:       post.Author.ID,
			Username: post.Author.Username,
		}
	}

	return resp
}
