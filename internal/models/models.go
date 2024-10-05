package models

type User struct {
	Id       int64
	Username string
	Password string
}

type Room struct {
	Uuid     string
	Name     string
	Password string
	OwnerId  int64
}

func (r *Room) IsPrivate() bool {
	return len(r.Password) > 0
}
