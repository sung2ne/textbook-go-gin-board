
type CommentResponse struct {
    ID              uint               `json:"id"`
    PostID          uint               `json:"post_id"`
    ParentID        *uint              `json:"parent_id,omitempty"`
    Author          *AuthorInfo        `json:"author"`
    Content         string             `json:"content"`
    ContentHTML     string             `json:"content_html"` // 멘션 하이라이트 포함
    CreatedAt       time.Time          `json:"created_at"`
    Replies         []*CommentResponse `json:"replies,omitempty"`
}

func ToCommentResponse(comment *domain.Comment) *CommentResponse {
    resp := &CommentResponse{
        ID:        comment.ID,
        PostID:    comment.PostID,
        ParentID:  comment.ParentID,
        CreatedAt: comment.CreatedAt,
    }

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
