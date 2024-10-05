package storage

import "errors"

var (
	UserNotFound   = errors.New("user not found")
	UsernameExists = errors.New("username is already taken")
	RoomNotFound   = errors.New("room not found")
)
