package email

import (
	"context"
)

// Email 이메일 데이터
type Email struct {
	To       []string
	Subject  string
	Body     string
	HTMLBody string
	From     string
	ReplyTo  string
}

// Sender 이메일 발송 인터페이스
type Sender interface {
	Send(ctx context.Context, email Email) error
}
