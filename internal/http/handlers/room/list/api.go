package roomlist

import (
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/types"
)

type Response struct {
	api.Response
	Data `json:"data"`
}

type Data struct {
	RoomsResponse []*types.RoomView `json:"rooms"`
	Size          int               `json:"size"`
}
