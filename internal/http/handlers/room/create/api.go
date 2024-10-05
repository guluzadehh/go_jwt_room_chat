package roomcreate

import (
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/types"
)

type Request struct {
	Name     string `json:"name" validate:"required,max=20"`
	Password string `json:"password"`
}

type Response struct {
	api.Response
	Data Data `json:"data"`
}

type Data struct {
	Room *types.RoomView `json:"room"`
}
