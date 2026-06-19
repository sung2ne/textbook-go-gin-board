package main

import (
    "context"
    "fmt"
    "time"

    "yourproject/batch"
)

func main() {
    items := make([]int, 1000)
    for i := range items {
        items[i] = i
    }

    processor := batch.NewProcessor(10, func(ctx context.Context, id int) error {
        time.Sleep(10 * time.Millisecond)
        return nil
    })

    start := time.Now()
    processor.Process(context.Background(), items)
    fmt.Printf("병렬 처리 (10 워커): %v\n", time.Since(start))
    // 약 1초 소요 (10배 빠름)
}
