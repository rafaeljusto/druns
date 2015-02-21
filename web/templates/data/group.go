package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/place"
)

type Group struct {
	Logged
	Form
	Group  group.Group
	Places []place.Place
}

func NewGroup(username core.Name, menu Menu) Group {
	return Group{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
