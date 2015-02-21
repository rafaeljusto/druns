package data

import (
	"github.com/rafaeljusto/druns/core"
)

type Schedule struct {
	Logged
}

func NewSchedule(username core.Name, menu Menu) Schedule {
	return Schedule{
		Logged: NewLogged(username, menu),
	}
}
