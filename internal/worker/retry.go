package worker

import (
    "context"
    "encoding/json"
    "log"
    "time"
)

const MaxRetries = 3

type RetryableHandler struct {
    handler    TaskHandler
    queue      Queue
    retryDelay time.Duration
}

func NewRetryableHandler(handler TaskHandler, queue Queue) *RetryableHandler {
    return &RetryableHandler{
        handler:    handler,
        queue:      queue,
        retryDelay: 1 * time.Minute,
    }
}

func (h *RetryableHandler) Handle(ctx context.Context, task *Task) error {
    err := h.handler(ctx, task.Payload)
    if err != nil {
        if task.Retries < MaxRetries {
            // 재시도 큐에 추가
            task.Retries++
            log.Printf("태스크 %s 재시도 예정 (%d/%d)", task.ID, task.Retries, MaxRetries)

            go func() {
                time.Sleep(h.retryDelay * time.Duration(task.Retries))
                h.queue.Enqueue(context.Background(), *task)
            }()
        } else {
            log.Printf("태스크 %s 최종 실패", task.ID)
            // Dead Letter Queue로 이동하거나 알림 발송
        }
        return err
    }
    return nil
}
