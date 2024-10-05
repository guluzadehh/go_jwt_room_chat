package login

import "github.com/guluzadehh/go_chat/internal/lib/api"

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	api.Response
	Data `json:"data"`
}

type Data struct {
	Token string `json:"token"`
}
