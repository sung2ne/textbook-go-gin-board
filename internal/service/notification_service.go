
type NotificationService interface {
    // ...
    CreateMentionNotifications(ctx context.Context, content string, postID uint, commentID uint, actorID uint) error
}

// CreateMentionNotifications 멘션 알림 생성
func (s *notificationService) CreateMentionNotifications(
    ctx context.Context,
    content string,
    postID uint,
    commentID uint,
    actorID uint,
) error {
    // 멘션된 사용자명 추출
    usernames := util.ParseMentions(content)
    if len(usernames) == 0 {
        return nil
    }

    // 사용자 조회
    users, err := s.userRepo.FindByUsernames(ctx, usernames)
    if err != nil {
        return err
    }

    // 행위자 정보 조회
    actor, err := s.userRepo.FindByID(ctx, actorID)
    if err != nil {
        return err
    }

    // 각 사용자에게 알림 생성
    for _, user := range users {
        // 본인을 멘션하면 알림 안 함
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

        // 에러가 발생해도 다른 사용자 알림은 계속 생성
        _ = s.notificationRepo.Create(ctx, notification)
    }

    return nil
}
