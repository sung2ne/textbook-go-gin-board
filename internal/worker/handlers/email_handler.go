package handlers

import (
    "context"
    "encoding/json"
    "time"

    "yourproject/internal/email"
)

type EmailPayload struct {
    To       []string `json:"to"`
    Subject  string   `json:"subject"`
    Body     string   `json:"body"`
    HTMLBody string   `json:"html_body,omitempty"`
}

type EmailHandler struct {
    sender email.Sender
}

func NewEmailHandler(sender email.Sender) *EmailHandler {
    return &EmailHandler{sender: sender}
}

func (h *EmailHandler) Handle(ctx context.Context, payload json.RawMessage) error {
    var data EmailPayload
    if err := json.Unmarshal(payload, &data); err != nil {
        return err
    }

    // 타임아웃 설정
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    return h.sender.Send(ctx, email.Email{
        To:       data.To,
        Subject:  data.Subject,
        Body:     data.Body,
        HTMLBody: data.HTMLBody,
    })
}
