package data

import (
	"github.com/rafaeljusto/druns/core/class"
	"github.com/rafaeljusto/druns/core/types"
)

type Schedule struct {
	Logged
	Classes []class.Class
}

func NewSchedule(username types.Name, menu Menu, classes []class.Class) Schedule {
	return Schedule{
		Logged:  NewLogged(username, menu),
		Classes: classes,
	}
}
