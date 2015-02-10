package data

import "github.com/rafaeljusto/druns/core/model"

type Schedule struct {
	Logged
}

func NewSchedule(username model.Name, menu Menu) Schedule {
	return Schedule{
		Logged: NewLogged(username, menu),
	}
}
