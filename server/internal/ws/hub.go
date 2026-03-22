package ws

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(msg *Message) {
	data, _ := json.Marshal(msg)
	h.broadcast <- data
}

func (h *Hub) BroadcastToType(msgType MessageType, data interface{}) {
	msg := &Message{Type: msgType}
	msg.Data, _ = json.Marshal(data)
	fullMsg, _ := json.Marshal(msg)
	h.mu.RLock()
	for client := range h.clients {
		if client.IsSubscribed(msgType) {
			select {
			case client.send <- fullMsg:
			default:
			}
		}
	}
	h.mu.RUnlock()
}

func (h *Hub) BroadcastToTypeRaw(msgType MessageType, data json.RawMessage) {
	h.mu.RLock()
	for client := range h.clients {
		if client.IsSubscribed(msgType) {
			fullMsg, _ := json.Marshal(&Message{Type: msgType, Data: data})
			select {
			case client.send <- fullMsg:
			default:
			}
		}
	}
	h.mu.RUnlock()
}
