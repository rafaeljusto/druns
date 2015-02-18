package data

import "github.com/rafaeljusto/druns/core/model"

type Group struct {
	Logged
	Form
	Group  model.Group
	Places []model.Place
}

func NewGroup(username model.Name, menu Menu) Group {
	return Group{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
