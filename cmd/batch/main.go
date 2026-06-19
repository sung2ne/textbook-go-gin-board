package main

import (
    "context"
    "fmt"
    "time"

    "yourproject/batch"
)

func main() {
    ctx := context.Background()
    items := make([]int, 1000)
    for i := range items {
        items[i] = i
    }

    progress := &batch.Progress{
        Total:     int64(len(items)),
        StartedAt: time.Now(),
    }

    // 진행 상황 출력 고루틴
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                fmt.Printf("\r진행: %.1f%% (%d/%d), 예상 남은 시간: %v",
                    progress.Percentage(),
                    progress.Completed,
                    progress.Total,
                    progress.ETA(),
                )
            }
        }
    }()

    processor := batch.NewProcessor(10, func(ctx context.Context, id int) error {
        time.Sleep(10 * time.Millisecond)
        progress.Increment()
        return nil
    })

    processor.Process(ctx, items)
    fmt.Println("\n완료!")
}
