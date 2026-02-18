package service

import (
	"context"

	"goboardapi/internal/auth"
	"goboardapi/internal/domain"
	"goboardapi/internal/dto"
	"goboardapi/internal/middleware"
	"goboardapi/internal/repository"
)

type CommentService interface {
	Create(ctx context.Context, postID uint, req *dto.CreateCommentRequest) (*domain.Comment, error)
	GetByPostID(ctx context.Context, postID uint) ([]*dto.CommentResponse, error)
	Update(ctx context.Context, id uint, req *dto.UpdateCommentRequest) (*domain.Comment, error)
	Delete(ctx context.Context, id uint) error
}

type commentService struct {
	commentRepo     repository.CommentRepository
	postRepo        repository.PostRepository
	notificationSvc NotificationService
}

func NewCommentService(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	notificationSvc NotificationService,
) CommentService {
	return &commentService{
		commentRepo:     commentRepo,
		postRepo:        postRepo,
		notificationSvc: notificationSvc,
	}
}

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

func (s *commentService) GetByPostID(ctx context.Context, postID uint) ([]*dto.CommentResponse, error) {
	_, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.FindByPostIDWithReplies(ctx, postID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.CommentResponse, len(comments))
	for i, comment := range comments {
		result[i] = dto.ToCommentResponse(&comment)
	}

	return result, nil
}

func (s *commentService) Update(ctx context.Context, id uint, req *dto.UpdateCommentRequest) (*domain.Comment, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 권한 검사: 본인만 수정 가능
	if comment.AuthorID != claims.UserID {
		return nil, ErrForbidden
	}

	comment.Content = req.Content

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) Delete(ctx context.Context, id uint) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	post, err := s.postRepo.FindByID(ctx, comment.PostID)
	if err != nil {
		return err
	}

	if !s.canDelete(comment, post, claims) {
		return ErrForbidden
	}

	// 대댓글이 있으면 소프트 삭제
	hasReplies, err := s.commentRepo.HasReplies(ctx, id)
	if err != nil {
		return err
	}

	if hasReplies {
		comment.IsDeleted = true
		comment.Content = ""
		return s.commentRepo.Update(ctx, comment)
	}

	return s.commentRepo.Delete(ctx, id)
}

// canDelete 삭제 권한 확인
func (s *commentService) canDelete(comment *domain.Comment, post *domain.Post, claims *auth.CustomClaims) bool {
	if comment.AuthorID == claims.UserID {
		return true
	}
	if post.AuthorID == claims.UserID {
		return true
	}
	if claims.Role == string(domain.RoleAdmin) {
		return true
	}
	return false
}
