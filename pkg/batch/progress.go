package batch

import (
	"sync/atomic"
	"time"
)

// Progress 진행 상황 추적
type Progress struct {
	Total     int64
	Completed int64
	Failed    int64
	StartedAt time.Time
}

// Increment 완료 카운터 증가
func (p *Progress) Increment() {
	atomic.AddInt64(&p.Completed, 1)
}

// IncrementFailed 실패 카운터 증가
func (p *Progress) IncrementFailed() {
	atomic.AddInt64(&p.Failed, 1)
}

// Percentage 진행률 반환
func (p *Progress) Percentage() float64 {
	if p.Total == 0 {
		return 0
	}
	return float64(atomic.LoadInt64(&p.Completed)) / float64(p.Total) * 100
}

// ETA 예상 남은 시간
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
