package email

import (
    "context"
)

type Email struct {
    To       []string
    Subject  string
    Body     string
    HTMLBody string
    From     string
    ReplyTo  string
}

type Sender interface {
    Send(ctx context.Context, email Email) error
}
