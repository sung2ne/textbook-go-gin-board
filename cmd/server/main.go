package main

import (
    "net/http"
    _ "net/http/pprof"
    "os"
)

func main() {
    // 개발 환경에서만 pprof 활성화
    if os.Getenv("ENV") == "development" {
        go func() {
            // localhost에서만 접근 가능
            http.ListenAndServe("localhost:6060", nil)
        }()
    }

    // 메인 서버
    // ...
}
