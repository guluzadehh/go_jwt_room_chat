package signup

import (
	"github.com/guluzadehh/go_chat/internal/lib/api"
	"github.com/guluzadehh/go_chat/internal/types"
)

type Request struct {
	Username     string `json:"username" validate:"required,max=16"`
	Password     string `json:"password" validate:"required,min=5,passwordpattern"`
	ConfPassword string `json:"conf_password" validate:"required,eqfield=Password"`
}

type Response struct {
	api.Response
	Data `json:"data"`
}

type Data struct {
	User *types.UserView `json:"user"`
}
