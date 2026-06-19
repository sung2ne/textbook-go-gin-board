package worker

import (
    "context"
    "sync"
)

type Queue interface {
    Enqueue(ctx context.Context, task Task) error
    Dequeue(ctx context.Context) (*Task, error)
    Close() error
}

type MemoryQueue struct {
    tasks chan Task
    done  chan struct{}
    once  sync.Once
}

func NewMemoryQueue(size int) *MemoryQueue {
    return &MemoryQueue{
        tasks: make(chan Task, size),
        done:  make(chan struct{}),
    }
}

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

func (q *MemoryQueue) Close() error {
    q.once.Do(func() {
        close(q.done)
    })
    return nil
}

var ErrQueueClosed = errors.New("queue is closed")
