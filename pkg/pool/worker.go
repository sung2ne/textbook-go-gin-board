package pool

import (
	"context"
	"sync"
)

// Job 작업 함수 타입
type Job func() error

// Result 작업 결과
type Result struct {
	Err error
}

// WorkerPool 워커 풀
type WorkerPool struct {
	numWorkers int
	jobs       chan Job
	results    chan Result
	wg         sync.WaitGroup
}

// NewWorkerPool 워커 풀 생성
func NewWorkerPool(numWorkers, queueSize int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan Job, queueSize),
		results:    make(chan Result, queueSize),
	}
}

// Start 워커 풀 시작
func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i+1)
	}
}

func (p *WorkerPool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			err := job()
			p.results <- Result{Err: err}
		}
	}
}

// Submit 작업 제출
func (p *WorkerPool) Submit(job Job) {
	p.jobs <- job
}

// Results 결과 채널 반환
func (p *WorkerPool) Results() <-chan Result {
	return p.results
}

// Close 작업 채널 닫기
func (p *WorkerPool) Close() {
	close(p.jobs)
}

// Wait 모든 워커 종료 대기
func (p *WorkerPool) Wait() {
	p.wg.Wait()
	close(p.results)
}
