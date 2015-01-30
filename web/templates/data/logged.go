package data

import "time"

const (
	MenuSchedule       Menu = "schedule"
	MenuReports        Menu = "reports"
	MenuGroups         Menu = "groups"
	MenuClients        Menu = "clients"
	MenuAdministrators Menu = "administrators"
)

type Menu string

type Logged struct {
	User string
	Menu string
	Time string
}

func NewLogged(user string, menu Menu) Logged {
	return Logged{
		User: user,
		Menu: string(menu),
		Time: time.Now().Format(time.RFC822),
	}
}
