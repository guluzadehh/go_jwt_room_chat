package roomchat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/guluzadehh/go_chat/internal/models"
	"github.com/guluzadehh/go_chat/internal/types"
)

type MessageType int

const JoinType MessageType = 0
const LeaveType MessageType = 1
const ClientType MessageType = 2

func (t *MessageType) String() string {
	switch *t {
	case JoinType:
		return "join"
	case LeaveType:
		return "leave"
	case ClientType:
		return "client"
	}

	return ""
}

func (t *MessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

type Message struct {
	Type      MessageType     `json:"type"`
	Msg       string          `json:"message"`
	From      *types.UserView `json:"from,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

func NewMessage(msg string, from *models.User) *Message {
	return &Message{
		Type:      ClientType,
		Msg:       msg,
		From:      types.NewUser(from),
		CreatedAt: time.Now(),
	}
}

func NewJoinMessage(u *models.User) *Message {
	return &Message{
		Type:      JoinType,
		Msg:       fmt.Sprintf("%s has joined the chat", u.Username),
		From:      nil,
		CreatedAt: time.Now(),
	}
}

func NewLeaveMessage(u *models.User) *Message {
	return &Message{
		Type:      LeaveType,
		Msg:       fmt.Sprintf("%s has left the chat", u.Username),
		From:      nil,
		CreatedAt: time.Now(),
	}
}
