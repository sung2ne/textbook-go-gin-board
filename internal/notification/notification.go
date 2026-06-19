package notification

import (
    "encoding/json"
    "time"
)

type NotificationType string

const (
    NotificationNewComment NotificationType = "new_comment"
    NotificationNewLike    NotificationType = "new_like"
    NotificationNewMessage NotificationType = "new_message"
    NotificationMention    NotificationType = "mention"
)

type Notification struct {
    ID        string           `json:"id"`
    Type      NotificationType `json:"type"`
    Title     string           `json:"title"`
    Body      string           `json:"body"`
    Data      interface{}      `json:"data,omitempty"`
    CreatedAt time.Time        `json:"created_at"`
}

func NewNotification(notifType NotificationType, title, body string, data interface{}) *Notification {
    return &Notification{
        ID:        uuid.New().String(),
        Type:      notifType,
        Title:     title,
        Body:      body,
        Data:      data,
        CreatedAt: time.Now(),
    }
}

func (n *Notification) JSON() []byte {
    data, _ := json.Marshal(n)
    return data
}
