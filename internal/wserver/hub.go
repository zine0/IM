package wserver

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type Hub struct {
	clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Send       chan *Message
	mu         sync.Mutex
}

type Message struct {
	From string
	To   string
	Time time.Time
	Msg  string
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Send:       make(chan *Message, 10),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case msg := <-h.Send:
			h.sendMessage(msg)
		}
	}
}

func (h *Hub) sendMessage(msg *Message) {
	to := msg.To

	client := h.clients[to]

	client.send <- msg

}

func (h *Hub) registerClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c.username] = c
	go c.RecvMessage()
	go c.SendMessage()
	zap.L().Info("client registered", zap.String("user", c.username))
}

func (h *Hub) unregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c.username)
	zap.L().Info("client unregistered", zap.String("user", c.username))
}
