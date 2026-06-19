
// PostListResponse에 하이라이트 필드 추가
type PostListResponse struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Author      string    `json:"author"`
    Views       int       `json:"views"`
    CreatedAt   time.Time `json:"created_at"`
    Highlight   string    `json:"highlight,omitempty"`  // 검색어 주변 텍스트
}
