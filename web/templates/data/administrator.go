package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/user"
)

type Administrator struct {
	Logged
	Form
	User user.User
}

func NewAdministrator(username core.Name, menu Menu) Administrator {
	return Administrator{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
