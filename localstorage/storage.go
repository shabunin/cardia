package localstorage

import "github.com/shabunin/cardia/user"

type Storage struct {
	user user.User
	fs   WriteFS
}

func NewLocalStorage(u user.User) *Storage {
	return &Storage{
		user: u,
	}
}
