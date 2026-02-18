package pool

import (
	"context"
	"sync"
)

// Task 제네릭 작업
type Task[T any, R any] struct {
	Input  T
	Result chan GenericResult[R]
}

// GenericResult 제네릭 결과
type GenericResult[R any] struct {
	Value R
	Err   error
}

// GenericPool 제네릭 워커 풀
type GenericPool[T any, R any] struct {
	numWorkers int
	tasks      chan Task[T, R]
	handler    func(context.Context, T) (R, error)
	wg         sync.WaitGroup
}

// NewGenericPool 제네릭 워커 풀 생성
func NewGenericPool[T any, R any](
	numWorkers int,
	queueSize int,
	handler func(context.Context, T) (R, error),
) *GenericPool[T, R] {
	return &GenericPool[T, R]{
		numWorkers: numWorkers,
		tasks:      make(chan Task[T, R], queueSize),
		handler:    handler,
	}
}

// Start 풀 시작
func (p *GenericPool[T, R]) Start(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
}

func (p *GenericPool[T, R]) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			result, err := p.handler(ctx, task.Input)
			task.Result <- GenericResult[R]{Value: result, Err: err}
			close(task.Result)
		}
	}
}

// Submit 작업 제출 (동기)
func (p *GenericPool[T, R]) Submit(ctx context.Context, input T) GenericResult[R] {
	resultCh := make(chan GenericResult[R], 1)
	p.tasks <- Task[T, R]{Input: input, Result: resultCh}

	select {
	case <-ctx.Done():
		return GenericResult[R]{Err: ctx.Err()}
	case r := <-resultCh:
		return r
	}
}

// Close 풀 종료
func (p *GenericPool[T, R]) Close() {
	close(p.tasks)
	p.wg.Wait()
}
