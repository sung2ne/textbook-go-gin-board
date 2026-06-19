package scheduler

import (
    "context"
    "log"
    "time"
)

type Job struct {
    Name     string
    Schedule time.Duration
    Handler  func(context.Context) error
}

type Scheduler struct {
    jobs []*Job
}

func New() *Scheduler {
    return &Scheduler{}
}

func (s *Scheduler) AddJob(job *Job) {
    s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Start(ctx context.Context) {
    for _, job := range s.jobs {
        go s.runJob(ctx, job)
    }
}

func (s *Scheduler) runJob(ctx context.Context, job *Job) {
    ticker := time.NewTicker(job.Schedule)
    defer ticker.Stop()

    // 시작 시 한 번 실행
    s.executeJob(ctx, job)

    for {
        select {
        case <-ctx.Done():
            log.Printf("스케줄러 종료: %s", job.Name)
            return
        case <-ticker.C:
            s.executeJob(ctx, job)
        }
    }
}

func (s *Scheduler) executeJob(ctx context.Context, job *Job) {
    log.Printf("작업 시작: %s", job.Name)
    start := time.Now()

    if err := job.Handler(ctx); err != nil {
        log.Printf("작업 실패: %s - %v", job.Name, err)
        return
    }

    log.Printf("작업 완료: %s (%v)", job.Name, time.Since(start))
}
