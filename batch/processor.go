package batch

import (
    "context"
    "sync"
)

type Processor[T any] struct {
    numWorkers int
    handler    func(context.Context, T) error
}

func NewProcessor[T any](numWorkers int, handler func(context.Context, T) error) *Processor[T] {
    return &Processor[T]{
        numWorkers: numWorkers,
        handler:    handler,
    }
}

func (p *Processor[T]) Process(ctx context.Context, items []T) error {
    jobs := make(chan T, len(items))
    errs := make(chan error, len(items))
    var wg sync.WaitGroup

    // 워커 시작
    for i := 0; i < p.numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range jobs {
                select {
                case <-ctx.Done():
                    return
                default:
                    if err := p.handler(ctx, item); err != nil {
                        errs <- err
                    }
                }
            }
        }()
    }

    // 작업 제출
    for _, item := range items {
        jobs <- item
    }
    close(jobs)

    // 완료 대기
    wg.Wait()
    close(errs)

    // 첫 번째 에러 반환
    for err := range errs {
        if err != nil {
            return err
        }
    }

    return nil
}
