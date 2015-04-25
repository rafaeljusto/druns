package data

import (
	"time"

	"github.com/rafaeljusto/druns/core/class"
	"github.com/rafaeljusto/druns/core/types"
)

type Schedule struct {
	Logged
	Begin    time.Time
	End      time.Time
	Classes  []class.Class
	Next     time.Time
	Previous time.Time
}

func NewSchedule(username types.Name, menu Menu, begin, end time.Time,
	classes []class.Class, next, previous time.Time) Schedule {

	return Schedule{
		Logged:   NewLogged(username, menu),
		Begin:    begin,
		End:      end,
		Classes:  classes,
		Next:     next,
		Previous: previous,
	}
}
