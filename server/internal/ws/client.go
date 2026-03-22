package ws

import (
	"sync"
)

type Client struct {
	hub       *Hub
	conn      interface{}
	send      chan []byte
	subscribe map[MessageType]bool
	mu        sync.RWMutex
}

func NewClient(hub *Hub) *Client {
	return &Client{
		hub:       hub,
		send:      make(chan []byte, 256),
		subscribe: make(map[MessageType]bool),
	}
}

func (c *Client) Subscribe(types []MessageType) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, t := range types {
		c.subscribe[t] = true
	}
}

func (c *Client) Unsubscribe(types []MessageType) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, t := range types {
		delete(c.subscribe, t)
	}
}

func (c *Client) IsSubscribed(t MessageType) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.subscribe[t]
}

func (c *Client) Close() {
	close(c.send)
}

func (c *Client) GetSend() <-chan []byte {
	return c.send
}
