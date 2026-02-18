package ws

import (
	"sync"
)

// Client WebSocket 클라이언트
type Client struct {
	ID     string
	UserID uint
	Send   chan []byte
}

// Hub WebSocket 허브
type Hub struct {
	clients     map[string]*Client
	userClients map[uint][]*Client
	register    chan *Client
	unregister  chan *Client
	broadcast   chan []byte
	mu          sync.RWMutex
}

// NewHub 허브 생성
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[string]*Client),
		userClients: make(map[uint][]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan []byte),
	}
}

// Run 허브 실행 (고루틴으로 실행)
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
		}
	}
}

// SendToUser 특정 사용자에게 메시지 전송
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

// Register 클라이언트 등록 채널 반환
func (h *Hub) Register() chan<- *Client {
	return h.register
}

// Unregister 클라이언트 등록 해제 채널 반환
func (h *Hub) Unregister() chan<- *Client {
	return h.unregister
}
