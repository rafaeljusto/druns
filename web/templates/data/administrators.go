package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/user"
)

type Administrators struct {
	Logged
	Users []user.User
}

func NewAdministrators(username core.Name, menu Menu, users []user.User) Administrators {
	return Administrators{
		Logged: NewLogged(username, menu),
		Users:  users,
	}
}
