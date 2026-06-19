package dto

import "time"

// CreatePostRequest 게시글 작성 요청
// @Description 게시글 작성에 필요한 정보
type CreatePostRequest struct {
    Title   string `json:"title" example:"새 게시글 제목" binding:"required,min=1,max=100"`
    Content string `json:"content" example:"게시글 본문 내용입니다." binding:"required,min=1"`
}

// UpdatePostRequest 게시글 수정 요청
// @Description 게시글 수정 정보 (부분 수정 가능)
type UpdatePostRequest struct {
    Title   *string `json:"title,omitempty" example:"수정된 제목" binding:"omitempty,min=1,max=100"`
    Content *string `json:"content,omitempty" example:"수정된 내용" binding:"omitempty,min=1"`
}

// PostResponse 게시글 응답
// @Description 게시글 상세 정보
type PostResponse struct {
    ID        uint      `json:"id" example:"1"`
    Title     string    `json:"title" example:"첫 번째 게시글"`
    Content   string    `json:"content" example:"게시글 내용입니다."`
    Author    string    `json:"author" example:"홍길동"`
    ViewCount int       `json:"viewCount" example:"100"`
    CreatedAt time.Time `json:"createdAt" example:"2024-01-15T10:30:00Z"`
    UpdatedAt time.Time `json:"updatedAt" example:"2024-01-15T10:30:00Z"`
}

// ListPostsResponse 게시글 목록 응답
// @Description 페이징된 게시글 목록
type ListPostsResponse struct {
    Data       []PostResponse `json:"data"`
    TotalCount int64          `json:"totalCount" example:"100"`
    Page       int            `json:"page" example:"1"`
    Size       int            `json:"size" example:"10"`
}

// ErrorResponse 에러 응답
// @Description API 에러 정보
type ErrorResponse struct {
    Error   string `json:"error" example:"잘못된 요청입니다"`
    Code    string `json:"code,omitempty" example:"INVALID_INPUT"`
    Details string `json:"details,omitempty" example:"title 필드는 필수입니다"`
}
