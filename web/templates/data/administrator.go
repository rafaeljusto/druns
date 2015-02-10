package data

import "github.com/rafaeljusto/druns/core/model"

type Administrator struct {
	Logged
	Form
	User model.User
}

func NewAdministrator(username model.Name, menu Menu) Administrator {
	return Administrator{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
