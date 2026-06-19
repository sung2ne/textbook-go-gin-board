package ws

import (
    "encoding/json"
    "log"
    "time"

    "github.com/google/uuid"
    "github.com/gorilla/websocket"
)

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512
)

func ServeWS(hub *Hub, conn *websocket.Conn, userID uint) {
    client := &Client{
        ID:     uuid.New().String(),
        UserID: userID,
        Conn:   conn,
        Send:   make(chan []byte, 256),
    }

    hub.register <- client

    go client.writePump()
    go client.readPump(hub)
}

func (c *Client) readPump(hub *Hub) {
    defer func() {
        hub.unregister <- c
        c.Conn.Close()
    }()

    c.Conn.SetReadLimit(maxMessageSize)
    c.Conn.SetReadDeadline(time.Now().Add(pongWait))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("에러: %v", err)
            }
            break
        }

        // 메시지 처리
        c.handleMessage(hub, message)
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }

        case <-ticker.C:
            c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Client) handleMessage(hub *Hub, message []byte) {
    // 메시지 타입에 따라 처리
    var msg map[string]interface{}
    if err := json.Unmarshal(message, &msg); err != nil {
        return
    }

    // 예: 채팅 메시지
    if msg["type"] == "chat" {
        hub.broadcast <- message
    }
}
