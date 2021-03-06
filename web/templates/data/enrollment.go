package data

import (
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/types"
)

type Enrollment struct {
	Logged
	Form
	Enrollment enrollment.Enrollment
	Clients    []client.Client
	Groups     []group.Group
	Back       string
}

func NewEnrollment(username types.Name, menu Menu) Enrollment {
	return Enrollment{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
