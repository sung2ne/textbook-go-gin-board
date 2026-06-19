package main

import (
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, Go!"))
    })

    log.Println("서버 시작: http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
