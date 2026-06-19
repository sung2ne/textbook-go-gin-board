package service

import (
    "context"

    "yourproject/internal/notification"
    "yourproject/internal/ws"
)

type NotificationService struct {
    hub        *ws.Hub
    notifRepo  repository.NotificationRepository
}

func NewNotificationService(hub *ws.Hub, repo repository.NotificationRepository) *NotificationService {
    return &NotificationService{
        hub:       hub,
        notifRepo: repo,
    }
}

func (s *NotificationService) NotifyNewComment(ctx context.Context, postAuthorID uint, comment *domain.Comment) error {
    notif := notification.NewNotification(
        notification.NotificationNewComment,
        "새 댓글",
        fmt.Sprintf("%s님이 댓글을 남겼습니다.", comment.Author.Username),
        map[string]interface{}{
            "post_id":    comment.PostID,
            "comment_id": comment.ID,
        },
    )

    // DB에 저장
    if err := s.notifRepo.Create(ctx, postAuthorID, notif); err != nil {
        return err
    }

    // 실시간 전송
    s.hub.SendToUser(postAuthorID, notif.JSON())

    return nil
}

func (s *NotificationService) NotifyNewLike(ctx context.Context, postAuthorID uint, likedBy *domain.User) error {
    notif := notification.NewNotification(
        notification.NotificationNewLike,
        "좋아요",
        fmt.Sprintf("%s님이 게시글에 좋아요를 눌렀습니다.", likedBy.Username),
        nil,
    )

    if err := s.notifRepo.Create(ctx, postAuthorID, notif); err != nil {
        return err
    }

    s.hub.SendToUser(postAuthorID, notif.JSON())

    return nil
}
