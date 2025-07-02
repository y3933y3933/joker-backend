package ws

import "sync"

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetRoom(code string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[code]
}

func (h *Hub) CreateRoom(code string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := NewRoom(code)
	h.rooms[code] = room

	go room.Run()
	return room

}

func (h *Hub) DeleteRoom(code string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, code)
}
