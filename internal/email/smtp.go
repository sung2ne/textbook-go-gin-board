package email

import (
    "context"
    "crypto/tls"
    "fmt"
    "net/smtp"
)

type SMTPConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
    UseTLS   bool
}

type SMTPSender struct {
    config SMTPConfig
}

func NewSMTPSender(config SMTPConfig) *SMTPSender {
    return &SMTPSender{config: config}
}

func (s *SMTPSender) Send(ctx context.Context, email Email) error {
    addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

    from := email.From
    if from == "" {
        from = s.config.From
    }

    // 이메일 본문 구성
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

    // context 타임아웃 처리
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
