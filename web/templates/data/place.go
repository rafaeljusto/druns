package data

import "github.com/rafaeljusto/druns/core/model"

type Place struct {
	Logged
	Form
	Place model.Place
}

func NewPlace(username model.Name, menu Menu) Place {
	return Place{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
