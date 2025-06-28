package ws

import (
	"encoding/json"
	"sync"
)

type Room struct {
	Code        string
	clients     map[*Client]bool
	clientsByID map[int64]*Client
	join        chan *Client
	leave       chan *Client
	broadcast   chan []byte
	mu          sync.RWMutex
}

func NewRoom(code string) *Room {
	return &Room{
		Code:        code,
		clients:     make(map[*Client]bool),
		clientsByID: make(map[int64]*Client),
		join:        make(chan *Client),
		leave:       make(chan *Client),
		broadcast:   make(chan []byte),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.join:
			r.mu.Lock()
			r.clients[client] = true
			r.clientsByID[client.ID] = client
			r.mu.Unlock()

		case client := <-r.leave:
			r.mu.Lock()
			delete(r.clients, client)
			delete(r.clientsByID, client.ID)
			r.mu.Unlock()

		case msg := <-r.broadcast:
			r.mu.RLock()
			for c := range r.clients {
				c.send <- msg
			}
			r.mu.RUnlock()
		}
	}
}

func (r *Room) Broadcast(msg any) {
	data, _ := json.Marshal(msg)
	r.broadcast <- data
}

func (r *Room) SendTo(playerID int64, msg any) {
	data, _ := json.Marshal(msg)
	r.mu.RLock()
	defer r.mu.RUnlock()
	if c, ok := r.clientsByID[playerID]; ok {
		c.send <- data
	}
}
