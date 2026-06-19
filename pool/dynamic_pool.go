package pool

import (
    "context"
    "sync"
    "sync/atomic"
)

type DynamicPool struct {
    minWorkers int
    maxWorkers int
    jobs       chan func()
    activeJobs int64
    workerCount int64
    mu         sync.Mutex
}

func NewDynamicPool(minWorkers, maxWorkers, queueSize int) *DynamicPool {
    p := &DynamicPool{
        minWorkers: minWorkers,
        maxWorkers: maxWorkers,
        jobs:       make(chan func(), queueSize),
    }

    // 최소 워커 시작
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

func (p *DynamicPool) Submit(job func()) {
    // 대기 작업이 많으면 워커 추가
    queueLen := len(p.jobs)
    workers := atomic.LoadInt64(&p.workerCount)

    if queueLen > int(workers) && int(workers) < p.maxWorkers {
        go p.worker()
        atomic.AddInt64(&p.workerCount, 1)
    }

    p.jobs <- job
}

func (p *DynamicPool) Stats() (workers, activeJobs, queueLen int64) {
    return atomic.LoadInt64(&p.workerCount),
           atomic.LoadInt64(&p.activeJobs),
           int64(len(p.jobs))
}
