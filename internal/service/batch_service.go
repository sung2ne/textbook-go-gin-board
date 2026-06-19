package service

import (
    "context"
    "log"

    "yourproject/internal/domain"
    "yourproject/internal/repository"
)

type BatchService struct {
    userRepo repository.UserRepository
}

func (s *BatchService) ProcessInactiveUsers(ctx context.Context) error {
    const batchSize = 100
    var lastID uint = 0

    for {
        // 배치 조회
        users, err := s.userRepo.FindInactiveBatch(ctx, lastID, batchSize)
        if err != nil {
            return err
        }

        if len(users) == 0 {
            break // 더 이상 데이터 없음
        }

        // 배치 처리
        for _, user := range users {
            if err := s.processInactiveUser(ctx, user); err != nil {
                log.Printf("처리 실패: %d - %v", user.ID, err)
            }
        }

        // 다음 배치를 위한 마지막 ID 업데이트
        lastID = users[len(users)-1].ID

        log.Printf("배치 처리 완료: lastID=%d", lastID)
    }

    return nil
}

func (s *BatchService) processInactiveUser(ctx context.Context, user *domain.User) error {
    // 비활성 사용자 처리 로직
    return nil
}
