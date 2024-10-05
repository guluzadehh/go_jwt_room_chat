package types

import (
	"github.com/guluzadehh/go_chat/internal/models"
)

type UserView struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

func NewUser(u *models.User) *UserView {
	if u == nil {
		return nil
	}

	return &UserView{
		Id:       u.Id,
		Username: u.Username,
	}
}

type RoomView struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	IsPrivate bool      `json:"is_private"`
	Owner     *UserView `json:"owner"`
}

func NewRoom(r *models.Room, o *models.User) *RoomView {
	if r == nil {
		return nil
	}

	return &RoomView{
		Uuid:      r.Uuid,
		Name:      r.Name,
		IsPrivate: r.IsPrivate(),
		Owner:     NewUser(o),
	}
}
