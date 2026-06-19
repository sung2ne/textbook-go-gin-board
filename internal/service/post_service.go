
func (s *postService) GetByID(ctx context.Context, id uint) (*dto.PostResponse, error) {
    post, err := s.postRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    resp := dto.ToPostResponse(post)

    // 로그인 상태면 좋아요 여부 확인
    if claims, ok := middleware.GetUserFromContext(ctx); ok {
        isLiked, _ := s.likeRepo.IsLiked(ctx, claims.UserID, id)
        resp.IsLiked = isLiked
    }

    return resp, nil
}
