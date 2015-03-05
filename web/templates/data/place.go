package data

import (
	"github.com/rafaeljusto/druns/core/place"
	"github.com/rafaeljusto/druns/core/types"
)

type Place struct {
	Logged
	Form
	Place place.Place
}

func NewPlace(username types.Name, menu Menu) Place {
	return Place{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
