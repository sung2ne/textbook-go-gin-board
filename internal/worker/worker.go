package worker

import (
    "context"
    "sync"
)

type Worker struct {
    wg     sync.WaitGroup
    ctx    context.Context
    cancel context.CancelFunc
}

func New() *Worker {
    ctx, cancel := context.WithCancel(context.Background())
    return &Worker{
        ctx:    ctx,
        cancel: cancel,
    }
}

func (w *Worker) Start(job func(ctx context.Context)) {
    w.wg.Add(1)
    go func() {
        defer w.wg.Done()
        job(w.ctx)
    }()
}

func (w *Worker) Shutdown() {
    w.cancel()       // 모든 작업에 취소 신호
    w.wg.Wait()      // 모든 작업 완료 대기
}
