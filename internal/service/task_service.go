package service

import (
    "context"
    "encoding/json"
    "time"

    "github.com/google/uuid"
    "yourproject/internal/worker"
)

type TaskService struct {
    queue worker.Queue
}

func NewTaskService(queue worker.Queue) *TaskService {
    return &TaskService{queue: queue}
}

// 이메일 발송 태스크 추가
func (s *TaskService) EnqueueEmail(ctx context.Context, to, subject, body string) error {
    payload, _ := json.Marshal(map[string]string{
        "to":      to,
        "subject": subject,
        "body":    body,
    })

    task := worker.Task{
        ID:        uuid.New().String(),
        Type:      worker.TaskSendEmail,
        Payload:   payload,
        CreatedAt: time.Now(),
    }

    return s.queue.Enqueue(ctx, task)
}
