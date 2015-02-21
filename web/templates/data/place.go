package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/place"
)

type Place struct {
	Logged
	Form
	Place place.Place
}

func NewPlace(username core.Name, menu Menu) Place {
	return Place{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
