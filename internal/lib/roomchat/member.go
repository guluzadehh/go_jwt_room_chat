package roomchat

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/guluzadehh/go_chat/internal/models"
)

type Member struct {
	room *ChatRoom

	conn     *websocket.Conn
	isClosed bool
	mu       sync.Mutex

	user *models.User
}

func NewMember(conn *websocket.Conn, user *models.User, room *ChatRoom) *Member {
	m := &Member{
		room:     room,
		conn:     conn,
		user:     user,
		isClosed: false,
	}

	ticker := time.NewTicker(room.hub.pingPeriod)
	go func() {
		defer func() {
			ticker.Stop()

			m.mu.Lock()
			if !m.isClosed {
				m.conn.Close()
				m.isClosed = true
			}
			m.mu.Unlock()
		}()

		for range ticker.C {
			m.conn.SetWriteDeadline(time.Now().Add(m.room.hub.writeWait))

			if err := m.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	return m
}

func (m *Member) WriteJSON(msg *Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isClosed {
		return
	}

	m.conn.SetWriteDeadline(time.Now().Add(m.room.hub.writeWait))
	if err := m.conn.WriteJSON(msg); err != nil {
		m.close()
	}
}

func (m *Member) ReadPump() {
	m.conn.SetReadDeadline(time.Now().Add(m.room.hub.pongWait))
	m.conn.SetPongHandler(func(string) error { m.conn.SetReadDeadline(time.Now().Add(m.room.hub.pongWait)); return nil })

	for {
		msgType, msg, err := m.conn.ReadMessage()
		if err != nil {
			m.mu.Lock()
			m.close()
			m.mu.Unlock()
			return
		}

		if msgType == websocket.CloseMessage {
			m.mu.Lock()
			if !m.isClosed {
				m.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				m.close()
			}
			m.mu.Unlock()
			return
		}

		m.room.Broadcast(m.NewMessage(string(msg)))
	}
}

func (m *Member) close() {
	if m.isClosed {
		return
	}

	m.conn.Close()
	m.isClosed = true
	m.room.Remove(m)
}

func (m *Member) NewMessage(rcv string) *Message {
	return NewMessage(rcv, m.user)
}
