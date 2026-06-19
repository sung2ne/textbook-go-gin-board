package dto

import "time"

// CreatePostRequest 게시글 작성 요청
type CreatePostRequest struct {
    // 게시글 제목 (1-100자)
    Title string `json:"title" example:"Go 언어 입문 가이드" binding:"required,min=1,max=100"`
    // 게시글 본문
    Content string `json:"content" example:"Go는 Google에서 개발한 프로그래밍 언어입니다. 간결한 문법과 강력한 동시성 지원이 특징입니다." binding:"required,min=1"`
}

// UpdatePostRequest 게시글 수정 요청
type UpdatePostRequest struct {
    // 수정할 제목 (선택)
    Title *string `json:"title,omitempty" example:"Go 언어 입문 가이드 (수정)"`
    // 수정할 본문 (선택)
    Content *string `json:"content,omitempty" example:"수정된 본문 내용입니다."`
}

// PostResponse 게시글 응답
type PostResponse struct {
    // 게시글 고유 ID
    ID uint `json:"id" example:"1"`
    // 게시글 제목
    Title string `json:"title" example:"Go 언어 입문 가이드"`
    // 게시글 본문
    Content string `json:"content" example:"Go는 Google에서 개발한 프로그래밍 언어입니다."`
    // 작성자 이름
    Author string `json:"author" example:"홍길동"`
    // 조회수
    ViewCount int `json:"viewCount" example:"152"`
    // 작성 일시
    CreatedAt time.Time `json:"createdAt" example:"2024-01-15T10:30:00Z"`
    // 수정 일시
    UpdatedAt time.Time `json:"updatedAt" example:"2024-01-15T14:20:00Z"`
}

// ListPostsResponse 게시글 목록 응답
type ListPostsResponse struct {
    // 게시글 목록
    Data []PostResponse `json:"data"`
    // 전체 게시글 수
    TotalCount int64 `json:"totalCount" example:"42"`
    // 현재 페이지 번호
    Page int `json:"page" example:"1"`
    // 페이지당 게시글 수
    Size int `json:"size" example:"10"`
}

// ErrorResponse 에러 응답
type ErrorResponse struct {
    // 에러 메시지
    Error string `json:"error" example:"게시글을 찾을 수 없습니다"`
    // 에러 코드
    Code string `json:"code,omitempty" example:"POST_NOT_FOUND"`
    // 상세 정보
    Details string `json:"details,omitempty" example:"요청한 ID: 999"`
}
