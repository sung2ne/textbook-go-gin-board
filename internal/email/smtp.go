package email

import (
	"context"
	"fmt"
	"net/smtp"
)

// SMTPConfig SMTP 설정
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

// SMTPSender SMTP 이메일 발송
type SMTPSender struct {
	config SMTPConfig
}

// NewSMTPSender SMTP 발송기 생성
func NewSMTPSender(config SMTPConfig) *SMTPSender {
	return &SMTPSender{config: config}
}

// Send 이메일 발송
func (s *SMTPSender) Send(ctx context.Context, email Email) error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	from := email.From
	if from == "" {
		from = s.config.From
	}

	headers := make(map[string]string)
	headers["From"] = from
	headers["Subject"] = email.Subject
	headers["MIME-Version"] = "1.0"

	var body string
	if email.HTMLBody != "" {
		headers["Content-Type"] = "text/html; charset=UTF-8"
		body = email.HTMLBody
	} else {
		headers["Content-Type"] = "text/plain; charset=UTF-8"
		body = email.Body
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(addr, auth, from, email.To, []byte(message))
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
