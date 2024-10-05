package refresh

import "github.com/guluzadehh/go_chat/internal/lib/api"

type Response struct {
	api.Response
	Data `json:"data"`
}

type Data struct {
	Token string `json:"token"`
}
