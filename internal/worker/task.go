package worker

import (
	"encoding/json"
	"time"
)

// TaskType 작업 유형
type TaskType string

const (
	TaskSendEmail     TaskType = "send_email"
	TaskSendPush      TaskType = "send_push"
	TaskGenerateThumb TaskType = "generate_thumbnail"
)

// Task 비동기 작업
type Task struct {
	ID        string          `json:"id"`
	Type      TaskType        `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
	Retries   int             `json:"retries"`
}
