package batch

import (
    "context"
    "sync/atomic"
    "time"
)

type Progress struct {
    Total     int64
    Completed int64
    Failed    int64
    StartedAt time.Time
}

func (p *Progress) Increment() {
    atomic.AddInt64(&p.Completed, 1)
}

func (p *Progress) IncrementFailed() {
    atomic.AddInt64(&p.Failed, 1)
}

func (p *Progress) Percentage() float64 {
    if p.Total == 0 {
        return 0
    }
    return float64(atomic.LoadInt64(&p.Completed)) / float64(p.Total) * 100
}

func (p *Progress) ETA() time.Duration {
    completed := atomic.LoadInt64(&p.Completed)
    if completed == 0 {
        return 0
    }

    elapsed := time.Since(p.StartedAt)
    remaining := p.Total - completed
    avgDuration := elapsed / time.Duration(completed)

    return avgDuration * time.Duration(remaining)
}
