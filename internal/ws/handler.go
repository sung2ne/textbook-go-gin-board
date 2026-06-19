package ws

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // 프로덕션에서는 적절한 검사 필요
    },
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket 업그레이드 실패: %v", err)
        return
    }
    defer conn.Close()

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("읽기 에러: %v", err)
            break
        }

        log.Printf("받은 메시지: %s", message)

        // 에코
        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Printf("쓰기 에러: %v", err)
            break
        }
    }
}
