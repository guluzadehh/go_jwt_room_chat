package roomchat

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/guluzadehh/go_chat/internal/models"
)

type ChatRoom struct {
	hub  *Hub
	room *models.Room

	members map[*Member]bool
	mu      sync.RWMutex

	cap int
}

func NewRoom(room *models.Room, hub *Hub) *ChatRoom {
	return &ChatRoom{
		hub:     hub,
		room:    room,
		members: make(map[*Member]bool),
		cap:     hub.cap,
	}
}

func (r *ChatRoom) Broadcast(msg *Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.broadcast(msg)
}

func (r *ChatRoom) NewMember(conn *websocket.Conn, user *models.User) (*Member, error) {
	m := NewMember(conn, user, r)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.add(m); err != nil {
		return nil, err
	}

	r.broadcast(NewJoinMessage(m.user))
	return m, nil
}

func (r *ChatRoom) Remove(m *Member) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.remove(m)
	r.broadcast(NewLeaveMessage(m.user))
}

func (r *ChatRoom) add(m *Member) error {
	if r.isFull() {
		return RoomIsFull
	}

	r.members[m] = true

	return nil
}

func (r *ChatRoom) remove(m *Member) {
	if _, ok := r.members[m]; !ok {
		return
	}

	delete(r.members, m)
	if r.isEmpty() {
		r.hub.DeleteRoom(r.room)
	}
}

func (r *ChatRoom) broadcast(msg *Message) {
	for m := range r.members {
		go m.WriteJSON(msg)
	}
}

func (r *ChatRoom) isFull() bool {
	return len(r.members) == r.cap
}

func (r *ChatRoom) isEmpty() bool {
	return len(r.members) == 0
}
