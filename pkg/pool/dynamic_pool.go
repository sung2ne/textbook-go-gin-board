package pool

import (
	"sync"
	"sync/atomic"
)

// DynamicPool 동적 워커 풀
type DynamicPool struct {
	minWorkers  int
	maxWorkers  int
	jobs        chan func()
	activeJobs  int64
	workerCount int64
	mu          sync.Mutex
}

// NewDynamicPool 동적 워커 풀 생성
func NewDynamicPool(minWorkers, maxWorkers, queueSize int) *DynamicPool {
	p := &DynamicPool{
		minWorkers: minWorkers,
		maxWorkers: maxWorkers,
		jobs:       make(chan func(), queueSize),
	}

	for i := 0; i < minWorkers; i++ {
		go p.worker()
		atomic.AddInt64(&p.workerCount, 1)
	}

	return p
}

func (p *DynamicPool) worker() {
	for job := range p.jobs {
		atomic.AddInt64(&p.activeJobs, 1)
		job()
		atomic.AddInt64(&p.activeJobs, -1)
	}
	atomic.AddInt64(&p.workerCount, -1)
}

// Submit 작업 제출
func (p *DynamicPool) Submit(job func()) {
	queueLen := len(p.jobs)
	workers := atomic.LoadInt64(&p.workerCount)

	if queueLen > int(workers) && int(workers) < p.maxWorkers {
		go p.worker()
		atomic.AddInt64(&p.workerCount, 1)
	}

	p.jobs <- job
}

// Stats 풀 상태 조회
func (p *DynamicPool) Stats() (workers, activeJobs, queueLen int64) {
	return atomic.LoadInt64(&p.workerCount),
		atomic.LoadInt64(&p.activeJobs),
		int64(len(p.jobs))
}
