package worker

import (
    "context"
    "encoding/json"
    "time"
)

type TaskType string

const (
    TaskSendEmail     TaskType = "send_email"
    TaskSendPush      TaskType = "send_push"
    TaskGenerateThumb TaskType = "generate_thumbnail"
)

type Task struct {
    ID        string          `json:"id"`
    Type      TaskType        `json:"type"`
    Payload   json.RawMessage `json:"payload"`
    CreatedAt time.Time       `json:"created_at"`
    Retries   int             `json:"retries"`
}
