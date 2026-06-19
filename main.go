package main

import (
    "fmt"
    "sync"
)

type Result struct {
    ID    int
    Value int
}

func main() {
    var wg sync.WaitGroup
    var mu sync.Mutex
    results := make([]Result, 0)

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            // 작업 수행
            value := id * 10

            // 결과 저장 (동기화 필요)
            mu.Lock()
            results = append(results, Result{ID: id, Value: value})
            mu.Unlock()
        }(i)
    }

    wg.Wait()

    fmt.Println("결과:")
    for _, r := range results {
        fmt.Printf("  ID: %d, Value: %d\n", r.ID, r.Value)
    }
}
