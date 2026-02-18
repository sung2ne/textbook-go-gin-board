package worker

import (
	"context"
	"errors"
	"sync"
)

// ErrQueueClosed 큐가 닫힌 경우
var ErrQueueClosed = errors.New("queue is closed")

// Queue 작업 큐 인터페이스
type Queue interface {
	Enqueue(ctx context.Context, task Task) error
	Dequeue(ctx context.Context) (*Task, error)
	Close() error
}

// MemoryQueue 메모리 기반 작업 큐
type MemoryQueue struct {
	tasks chan Task
	done  chan struct{}
	once  sync.Once
}

// NewMemoryQueue 메모리 큐 생성
func NewMemoryQueue(size int) *MemoryQueue {
	return &MemoryQueue{
		tasks: make(chan Task, size),
		done:  make(chan struct{}),
	}
}

// Enqueue 작업 추가
func (q *MemoryQueue) Enqueue(ctx context.Context, task Task) error {
	select {
	case q.tasks <- task:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-q.done:
		return ErrQueueClosed
	}
}

// Dequeue 작업 꺼내기
func (q *MemoryQueue) Dequeue(ctx context.Context) (*Task, error) {
	select {
	case task := <-q.tasks:
		return &task, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-q.done:
		return nil, ErrQueueClosed
	}
}

// Close 큐 닫기
func (q *MemoryQueue) Close() error {
	q.once.Do(func() {
		close(q.done)
	})
	return nil
}
