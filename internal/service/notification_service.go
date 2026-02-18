package service

import (
	"context"
	"fmt"

	"goboardapi/internal/domain"
	"goboardapi/internal/middleware"
	"goboardapi/internal/repository"
	"goboardapi/internal/util"
)

type NotificationService interface {
	CreateCommentNotification(ctx context.Context, post *domain.Post, comment *domain.Comment, actorID uint) error
	CreateReplyNotification(ctx context.Context, parentComment *domain.Comment, reply *domain.Comment, actorID uint) error
	CreateMentionNotifications(ctx context.Context, content string, postID uint, commentID uint, actorID uint) error

	GetNotifications(ctx context.Context, page, size int) ([]*domain.Notification, int64, error)
	GetUnreadCount(ctx context.Context) (int64, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

func (s *notificationService) CreateCommentNotification(
	ctx context.Context,
	post *domain.Post,
	comment *domain.Comment,
	actorID uint,
) error {
	if post.AuthorID == actorID {
		return nil
	}

	actor, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return err
	}

	notification := &domain.Notification{
		UserID:    post.AuthorID,
		Type:      domain.NotificationComment,
		Message:   fmt.Sprintf("%s님이 회원님의 게시글에 댓글을 남겼습니다.", actor.Username),
		Link:      fmt.Sprintf("/posts/%d", post.ID),
		ActorID:   actorID,
		PostID:    &post.ID,
		CommentID: &comment.ID,
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationService) CreateReplyNotification(
	ctx context.Context,
	parentComment *domain.Comment,
	reply *domain.Comment,
	actorID uint,
) error {
	if parentComment.AuthorID == actorID {
		return nil
	}

	actor, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return err
	}

	notification := &domain.Notification{
		UserID:    parentComment.AuthorID,
		Type:      domain.NotificationReply,
		Message:   fmt.Sprintf("%s님이 회원님의 댓글에 답글을 남겼습니다.", actor.Username),
		Link:      fmt.Sprintf("/posts/%d#comment-%d", parentComment.PostID, reply.ID),
		ActorID:   actorID,
		PostID:    &parentComment.PostID,
		CommentID: &reply.ID,
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationService) CreateMentionNotifications(
	ctx context.Context,
	content string,
	postID uint,
	commentID uint,
	actorID uint,
) error {
	usernames := util.ParseMentions(content)
	if len(usernames) == 0 {
		return nil
	}

	users, err := s.userRepo.FindByUsernames(ctx, usernames)
	if err != nil {
		return err
	}

	actor, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.ID == actorID {
			continue
		}

		notification := &domain.Notification{
			UserID:    user.ID,
			Type:      domain.NotificationMention,
			Message:   fmt.Sprintf("%s님이 회원님을 멘션했습니다.", actor.Username),
			Link:      fmt.Sprintf("/posts/%d#comment-%d", postID, commentID),
			ActorID:   actorID,
			PostID:    &postID,
			CommentID: &commentID,
		}

		_ = s.notificationRepo.Create(ctx, notification)
	}

	return nil
}

func (s *notificationService) GetNotifications(ctx context.Context, page, size int) ([]*domain.Notification, int64, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return nil, 0, ErrUnauthorized
	}

	offset := (page - 1) * size
	return s.notificationRepo.FindByUserID(ctx, claims.UserID, offset, size)
}

func (s *notificationService) GetUnreadCount(ctx context.Context) (int64, error) {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return 0, ErrUnauthorized
	}

	return s.notificationRepo.CountUnread(ctx, claims.UserID)
}

func (s *notificationService) MarkAsRead(ctx context.Context, id uint) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	return s.notificationRepo.MarkAsRead(ctx, id, claims.UserID)
}

func (s *notificationService) MarkAllAsRead(ctx context.Context) error {
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		return ErrUnauthorized
	}

	return s.notificationRepo.MarkAllAsRead(ctx, claims.UserID)
}
