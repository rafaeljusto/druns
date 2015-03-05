package data

import (
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/place"
	"github.com/rafaeljusto/druns/core/types"
)

type Group struct {
	Logged
	Form
	Group       group.Group
	Places      []place.Place
	Enrollments []enrollment.Enrollment
}

func NewGroup(username types.Name, menu Menu) Group {
	return Group{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
