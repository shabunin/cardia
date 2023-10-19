package service

import (
	"github.com/shabunin/cardia/user"
	"modernc.org/sqlite"
)

type Service interface {
	Name() string        // return service name
	Account() *user.User // return authenticated service account
	DB() sqlite.Driver   // return database for given service

	Start() error
	Stop() error
	Restart() error
}
