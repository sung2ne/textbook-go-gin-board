
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

    // 탈퇴한 사용자 처리
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
