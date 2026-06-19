package ws

import (
    "sync"
)

type Client struct {
    ID     string
    UserID uint
    Conn   *websocket.Conn
    Send   chan []byte
}

type Hub struct {
    clients    map[string]*Client
    userClients map[uint][]*Client // 사용자별 클라이언트
    register   chan *Client
    unregister chan *Client
    broadcast  chan []byte
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:     make(map[string]*Client),
        userClients: make(map[uint][]*Client),
        register:    make(chan *Client),
        unregister:  make(chan *Client),
        broadcast:   make(chan []byte),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.addClient(client)

        case client := <-h.unregister:
            h.removeClient(client)

        case message := <-h.broadcast:
            h.broadcastMessage(message)
        }
    }
}

func (h *Hub) addClient(client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    h.clients[client.ID] = client
    h.userClients[client.UserID] = append(h.userClients[client.UserID], client)
}

func (h *Hub) removeClient(client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if _, ok := h.clients[client.ID]; ok {
        delete(h.clients, client.ID)
        close(client.Send)

        // 사용자 클라이언트 목록에서도 제거
        clients := h.userClients[client.UserID]
        for i, c := range clients {
            if c.ID == client.ID {
                h.userClients[client.UserID] = append(clients[:i], clients[i+1:]...)
                break
            }
        }
    }
}

func (h *Hub) broadcastMessage(message []byte) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    for _, client := range h.clients {
        select {
        case client.Send <- message:
        default:
            // 버퍼가 가득 찬 경우
        }
    }
}

// 특정 사용자에게 메시지 전송
func (h *Hub) SendToUser(userID uint, message []byte) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    clients := h.userClients[userID]
    for _, client := range clients {
        select {
        case client.Send <- message:
        default:
        }
    }
}
