package dto

import (
	"time"

	"goboardapi/internal/domain"
	"goboardapi/internal/util"
)

// CreateCommentRequest 댓글 생성 요청
type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

// UpdateCommentRequest 댓글 수정 요청
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CommentResponse 댓글 응답
type CommentResponse struct {
	ID          uint               `json:"id"`
	PostID      uint               `json:"post_id"`
	ParentID    *uint              `json:"parent_id,omitempty"`
	Author      *AuthorInfo        `json:"author"`
	Content     string             `json:"content"`
	ContentHTML string             `json:"content_html"`
	CreatedAt   time.Time          `json:"created_at"`
	Replies     []*CommentResponse `json:"replies,omitempty"`
}

// ToCommentResponse 도메인 → 응답 DTO 변환
func ToCommentResponse(comment *domain.Comment) *CommentResponse {
	resp := &CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
	}

	// 삭제된 댓글 처리
	if comment.IsDeleted {
		resp.Content = "삭제된 댓글입니다"
		resp.ContentHTML = resp.Content
		resp.Author = nil
	} else {
		resp.Content = comment.Content
		resp.ContentHTML = util.HighlightMentions(comment.Content)
		if comment.Author != nil {
			resp.Author = &AuthorInfo{
				ID:       comment.Author.ID,
				Username: comment.Author.Username,
			}
		}
	}

	// 대댓글 변환
	if len(comment.Replies) > 0 {
		for _, reply := range comment.Replies {
			resp.Replies = append(resp.Replies, ToCommentResponse(&reply))
		}
	}

	return resp
}
