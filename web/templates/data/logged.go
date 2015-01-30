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
	Username string
	Menu     string
	Time     string
}

func NewLogged(username string, menu Menu) Logged {
	return Logged{
		Username: username,
		Menu:     string(menu),
		Time:     time.Now().Format(time.RFC822),
	}
}
