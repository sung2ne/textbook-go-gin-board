package email

import (
    "context"
    "log"
)

type ConsoleSender struct{}

func NewConsoleSender() *ConsoleSender {
    return &ConsoleSender{}
}

func (s *ConsoleSender) Send(ctx context.Context, email Email) error {
    log.Printf(`
========== 이메일 발송 ==========
To: %v
Subject: %s
Body: %s
================================
`, email.To, email.Subject, email.Body)
    return nil
}
