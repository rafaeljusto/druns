package data

import (
	"time"

	"github.com/rafaeljusto/druns/core"
)

const (
	MenuSchedule       Menu = "schedule"
	MenuReports        Menu = "reports"
	MenuGroups         Menu = "groups"
	MenuClients        Menu = "clients"
	MenuPlaces         Menu = "places"
	MenuAdministrators Menu = "administrators"
)

type Menu string

type Logged struct {
	Username string
	Menu     string
	Time     string
}

func NewLogged(username core.Name, menu Menu) Logged {
	return Logged{
		Username: username.String(),
		Menu:     string(menu),
		Time:     time.Now().Format(time.RFC822),
	}
}
