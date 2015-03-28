package data

import (
	"github.com/rafaeljusto/druns/core/class"
	"github.com/rafaeljusto/druns/core/types"
)

type Class struct {
	Logged
	Form
	Class class.Class
}

func NewClass(username types.Name, menu Menu) Class {
	return Class{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
