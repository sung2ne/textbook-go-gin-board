
func (s *commentService) Create(ctx context.Context, postID uint, req *dto.CreateCommentRequest) (*domain.Comment, error) {
    claims, ok := middleware.GetUserFromContext(ctx)
    if !ok {
        return nil, ErrUnauthorized
    }

    post, err := s.postRepo.FindByID(ctx, postID)
    if err != nil {
        return nil, err
    }

    var parentComment *domain.Comment
    if req.ParentID != nil {
        parentComment, err = s.commentRepo.FindByID(ctx, *req.ParentID)
        if err != nil {
            return nil, err
        }
    }

    comment := &domain.Comment{
        PostID:   postID,
        ParentID: req.ParentID,
        AuthorID: claims.UserID,
        Content:  req.Content,
    }

    if err := s.commentRepo.Create(ctx, comment); err != nil {
        return nil, err
    }

    // 댓글/대댓글 알림
    if req.ParentID != nil {
        _ = s.notificationSvc.CreateReplyNotification(ctx, parentComment, comment, claims.UserID)
    } else {
        _ = s.notificationSvc.CreateCommentNotification(ctx, post, comment, claims.UserID)
    }

    // 멘션 알림
    _ = s.notificationSvc.CreateMentionNotifications(ctx, req.Content, postID, comment.ID, claims.UserID)

    return s.commentRepo.FindByID(ctx, comment.ID)
}
