package worker

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

// TaskHandler 작업 핸들러 함수
type TaskHandler func(ctx context.Context, payload json.RawMessage) error

// Worker 작업 처리 워커
type Worker struct {
	queue      Queue
	handlers   map[TaskType]TaskHandler
	numWorkers int
	wg         sync.WaitGroup
}

// NewWorker 워커 생성
func NewWorker(queue Queue, numWorkers int) *Worker {
	return &Worker{
		queue:      queue,
		handlers:   make(map[TaskType]TaskHandler),
		numWorkers: numWorkers,
	}
}

// RegisterHandler 핸들러 등록
func (w *Worker) RegisterHandler(taskType TaskType, handler TaskHandler) {
	w.handlers[taskType] = handler
}

// Start 워커 시작
func (w *Worker) Start(ctx context.Context) {
	for i := 0; i < w.numWorkers; i++ {
		w.wg.Add(1)
		go w.run(ctx, i+1)
	}
}

func (w *Worker) run(ctx context.Context, id int) {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("워커 %d 종료", id)
			return
		default:
		}

		task, err := w.queue.Dequeue(ctx)
		if err != nil {
			if err == ErrQueueClosed || err == context.Canceled {
				return
			}
			continue
		}

		w.processTask(ctx, id, task)
	}
}

func (w *Worker) processTask(ctx context.Context, workerID int, task *Task) {
	handler, ok := w.handlers[task.Type]
	if !ok {
		log.Printf("워커 %d: 알 수 없는 태스크 타입: %s", workerID, task.Type)
		return
	}

	log.Printf("워커 %d: 태스크 %s 처리 시작", workerID, task.ID)

	if err := handler(ctx, task.Payload); err != nil {
		log.Printf("워커 %d: 태스크 %s 실패: %v", workerID, task.ID, err)
		return
	}

	log.Printf("워커 %d: 태스크 %s 완료", workerID, task.ID)
}

// Wait 모든 워커 종료 대기
func (w *Worker) Wait() {
	w.wg.Wait()
}
