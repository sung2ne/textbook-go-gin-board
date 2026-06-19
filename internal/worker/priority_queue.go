package worker

import (
    "container/heap"
    "context"
    "sync"
)

type Priority int

const (
    PriorityLow    Priority = 0
    PriorityNormal Priority = 1
    PriorityHigh   Priority = 2
)

type PriorityTask struct {
    Task     Task
    Priority Priority
    Index    int
}

type PriorityQueue struct {
    items []*PriorityTask
    mu    sync.Mutex
    cond  *sync.Cond
}

func NewPriorityQueue() *PriorityQueue {
    pq := &PriorityQueue{
        items: make([]*PriorityTask, 0),
    }
    pq.cond = sync.NewCond(&pq.mu)
    heap.Init(pq)
    return pq
}

func (pq *PriorityQueue) Len() int { return len(pq.items) }

func (pq *PriorityQueue) Less(i, j int) bool {
    return pq.items[i].Priority > pq.items[j].Priority // 높은 우선순위 먼저
}

func (pq *PriorityQueue) Swap(i, j int) {
    pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
    pq.items[i].Index = i
    pq.items[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
    item := x.(*PriorityTask)
    item.Index = len(pq.items)
    pq.items = append(pq.items, item)
}

func (pq *PriorityQueue) Pop() any {
    old := pq.items
    n := len(old)
    item := old[n-1]
    old[n-1] = nil
    pq.items = old[0 : n-1]
    return item
}

func (pq *PriorityQueue) Enqueue(priority Priority, task Task) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    heap.Push(pq, &PriorityTask{Task: task, Priority: priority})
    pq.cond.Signal()
}

func (pq *PriorityQueue) Dequeue(ctx context.Context) (*Task, error) {
    pq.mu.Lock()
    defer pq.mu.Unlock()

    for len(pq.items) == 0 {
        pq.cond.Wait()
    }

    item := heap.Pop(pq).(*PriorityTask)
    return &item.Task, nil
}
