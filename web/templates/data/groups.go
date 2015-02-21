package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/group"
)

type Groups struct {
	Logged
	Groups []group.Group
}

func NewGroups(username core.Name, menu Menu, groups []group.Group) Groups {
	return Groups{
		Logged: NewLogged(username, menu),
		Groups: groups,
	}
}
