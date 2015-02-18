package data

import "github.com/rafaeljusto/druns/core/model"

type Places struct {
	Logged
	Places []model.Place
}

func NewPlaces(username model.Name, menu Menu, places []model.Place) Places {
	return Places{
		Logged: NewLogged(username, menu),
		Places: places,
	}
}
