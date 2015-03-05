package data

import (
	"github.com/rafaeljusto/druns/core/types"
	"github.com/rafaeljusto/druns/core/user"
)

type Administrators struct {
	Logged
	Users []user.User
}

func NewAdministrators(username types.Name, menu Menu, users []user.User) Administrators {
	return Administrators{
		Logged: NewLogged(username, menu),
		Users:  users,
	}
}
