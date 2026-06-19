package dto

// ErrorResponse 에러 응답
type ErrorResponse struct {
    Error   string `json:"error" example:"잘못된 요청입니다"`
    Code    string `json:"code,omitempty" example:"INVALID_INPUT"`
    Details string `json:"details,omitempty" example:"title 필드는 필수입니다"`
}

// ListResponse 목록 응답
type ListResponse[T any] struct {
    Data       []T   `json:"data"`
    TotalCount int64 `json:"totalCount" example:"100"`
    Page       int   `json:"page" example:"1"`
    Size       int   `json:"size" example:"10"`
}
