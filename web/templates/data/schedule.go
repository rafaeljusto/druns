package data

type Schedule struct {
	Logged
}

func NewSchedule(user string, menu Menu) Schedule {
	return Schedule{
		Logged: NewLogged(user, menu),
	}
}
