package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "yourproject/internal/ws"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func (h *WSHandler) HandleWebSocket(c *gin.Context) {
    // 인증 확인
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "인증 필요"})
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }

    ws.ServeWS(h.hub, conn, userID.(uint))
}
