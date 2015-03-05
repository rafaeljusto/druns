package data

import (
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/types"
)

type Groups struct {
	Logged
	Groups []group.Group
}

func NewGroups(username types.Name, menu Menu, groups []group.Group) Groups {
	return Groups{
		Logged: NewLogged(username, menu),
		Groups: groups,
	}
}
