package roomchat

import (
	"sync"
	"time"

	"github.com/guluzadehh/go_chat/internal/config"
	"github.com/guluzadehh/go_chat/internal/models"
)

type Hub struct {
	rooms map[string]*ChatRoom
	mu    sync.RWMutex

	cap int

	writeWait  time.Duration
	pongWait   time.Duration
	pingPeriod time.Duration
}

func NewHub(config *config.Config) *Hub {
	return &Hub{
		rooms:      make(map[string]*ChatRoom),
		cap:        config.Chat.Room.Capacity,
		writeWait:  config.Chat.WriteWait,
		pongWait:   config.Chat.PongWait,
		pingPeriod: config.Chat.PingPeriod,
	}
}

func (h *Hub) GetOrCreateRoom(r *models.Room) *ChatRoom {
	h.mu.RLock()
	room, ok := h.rooms[r.Uuid]
	h.mu.RUnlock()

	if ok {
		return room
	}

	room = NewRoom(r, h)

	h.mu.Lock()
	h.rooms[r.Uuid] = room
	h.mu.Unlock()

	return room
}

func (h *Hub) DeleteRoom(r *models.Room) {
	h.mu.Lock()
	delete(h.rooms, r.Uuid)
	h.mu.Unlock()
}
