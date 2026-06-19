package main

import (
    "os"
    "runtime/trace"
)

func main() {
    // 트레이스 파일 생성
    f, err := os.Create("trace.out")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    // 트레이스 시작
    if err := trace.Start(f); err != nil {
        panic(err)
    }
    defer trace.Stop()

    // 측정할 코드
    doWork()
}

func doWork() {
    // 작업 수행
}
