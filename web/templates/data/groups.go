package data

import "github.com/rafaeljusto/druns/core/model"

type Groups struct {
	Logged
	Groups []model.Group
}

func NewGroups(username model.Name, menu Menu, groups []model.Group) Groups {
	return Groups{
		Logged: NewLogged(username, menu),
		Groups: groups,
	}
}
