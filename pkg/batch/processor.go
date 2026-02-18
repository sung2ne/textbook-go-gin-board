package batch

import (
	"context"
	"sync"
)

// Processor 병렬 배치 프로세서
type Processor[T any] struct {
	numWorkers int
	handler    func(context.Context, T) error
}

// NewProcessor 프로세서 생성
func NewProcessor[T any](numWorkers int, handler func(context.Context, T) error) *Processor[T] {
	return &Processor[T]{
		numWorkers: numWorkers,
		handler:    handler,
	}
}

// Process 배치 처리
func (p *Processor[T]) Process(ctx context.Context, items []T) error {
	jobs := make(chan T, len(items))
	errs := make(chan error, len(items))
	var wg sync.WaitGroup

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

	for _, item := range items {
		jobs <- item
	}
	close(jobs)

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
