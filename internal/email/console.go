package email

import (
	"context"
	"log"
)

// ConsoleSender 콘솔 출력 이메일 발송 (개발용)
type ConsoleSender struct{}

// NewConsoleSender 콘솔 발송기 생성
func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

// Send 콘솔에 이메일 출력
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
