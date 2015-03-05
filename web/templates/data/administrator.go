package data

import (
	"github.com/rafaeljusto/druns/core/types"
	"github.com/rafaeljusto/druns/core/user"
)

type Administrator struct {
	Logged
	Form
	User user.User
}

func NewAdministrator(username types.Name, menu Menu) Administrator {
	return Administrator{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
