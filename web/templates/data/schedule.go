package data

import (
	"github.com/rafaeljusto/druns/core/types"
)

type Schedule struct {
	Logged
}

func NewSchedule(username types.Name, menu Menu) Schedule {
	return Schedule{
		Logged: NewLogged(username, menu),
	}
}
