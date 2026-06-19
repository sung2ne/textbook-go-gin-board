
// PostResponse 게시글 응답
type PostResponse struct {
    ID        uint        `json:"id"`
    Title     string      `json:"title"`
    Content   string      `json:"content"`
    Author    *AuthorInfo `json:"author"`
    Views     int         `json:"views"`
    LikeCount int         `json:"like_count"`
    IsLiked   bool        `json:"is_liked"` // 현재 사용자의 좋아요 여부
    CreatedAt time.Time   `json:"created_at"`
    UpdatedAt time.Time   `json:"updated_at"`
}
