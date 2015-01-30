package data

type Schedule struct {
	Logged
}

func NewSchedule(username string, menu Menu) Schedule {
	return Schedule{
		Logged: NewLogged(username, menu),
	}
}
