
const MaxReplyDepth = 3  // 최대 3단계까지

// Create 댓글 생성
func (s *CommentService) Create(postID uint, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
    // 게시글 존재 확인
    _, err := s.postRepo.FindByID(postID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrPostNotExists
        }
        return nil, err
    }

    // 부모 댓글 확인 (대댓글인 경우)
    if req.ParentID != nil {
        parent, err := s.commentRepo.FindByID(*req.ParentID)
        if err != nil {
            return nil, errors.New("부모 댓글을 찾을 수 없습니다")
        }
        if parent.PostID != postID {
            return nil, errors.New("부모 댓글이 다른 게시글에 속합니다")
        }

        // 깊이 확인
        depth := s.getCommentDepth(parent)
        if depth >= MaxReplyDepth {
            return nil, errors.New("더 이상 대댓글을 작성할 수 없습니다")
        }
    }

    comment := &domain.Comment{
        PostID:   postID,
        ParentID: req.ParentID,
        Content:  req.Content,
        Author:   req.Author,
    }

    if err := s.commentRepo.Create(comment); err != nil {
        return nil, err
    }

    return s.toResponse(comment), nil
}

// getCommentDepth 댓글 깊이 계산
func (s *CommentService) getCommentDepth(comment *domain.Comment) int {
    depth := 1
    current := comment

    for current.ParentID != nil {
        parent, err := s.commentRepo.FindByID(*current.ParentID)
        if err != nil {
            break
        }
        current = parent
        depth++
    }

    return depth
}
