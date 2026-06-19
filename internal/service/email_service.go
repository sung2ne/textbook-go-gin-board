package service

import (
    "context"
    "encoding/json"

    "yourproject/internal/email"
    "yourproject/internal/worker"
)

type EmailService struct {
    queue   worker.Queue
    appName string
}

func NewEmailService(queue worker.Queue, appName string) *EmailService {
    return &EmailService{
        queue:   queue,
        appName: appName,
    }
}

func (s *EmailService) SendWelcome(ctx context.Context, to, username string) error {
    html, err := email.RenderTemplate("welcome", email.TemplateData{
        Username: username,
        AppName:  s.appName,
    })
    if err != nil {
        return err
    }

    payload, _ := json.Marshal(handlers.EmailPayload{
        To:       []string{to},
        Subject:  s.appName + "에 오신 것을 환영합니다!",
        HTMLBody: html,
    })

    return s.queue.Enqueue(ctx, worker.Task{
        ID:      uuid.New().String(),
        Type:    worker.TaskSendEmail,
        Payload: payload,
    })
}

func (s *EmailService) SendPasswordReset(ctx context.Context, to, username, resetLink string) error {
    html, err := email.RenderTemplate("reset_password", email.TemplateData{
        Username: username,
        Link:     resetLink,
        AppName:  s.appName,
    })
    if err != nil {
        return err
    }

    payload, _ := json.Marshal(handlers.EmailPayload{
        To:       []string{to},
        Subject:  "[" + s.appName + "] 비밀번호 재설정",
        HTMLBody: html,
    })

    return s.queue.Enqueue(ctx, worker.Task{
        ID:      uuid.New().String(),
        Type:    worker.TaskSendEmail,
        Payload: payload,
    })
}
