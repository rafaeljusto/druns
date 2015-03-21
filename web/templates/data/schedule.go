package data

import (
	"time"

	"github.com/rafaeljusto/druns/core/class"
	"github.com/rafaeljusto/druns/core/types"
)

type Schedule struct {
	Logged
	Begin   time.Time
	End     time.Time
	Classes []class.Class
}

func NewSchedule(username types.Name, menu Menu, begin, end time.Time, classes []class.Class) Schedule {
	return Schedule{
		Logged:  NewLogged(username, menu),
		Begin:   begin,
		End:     end,
		Classes: classes,
	}
}
