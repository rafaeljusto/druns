package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/place"
)

type Places struct {
	Logged
	Places []place.Place
}

func NewPlaces(username core.Name, menu Menu, places []place.Place) Places {
	return Places{
		Logged: NewLogged(username, menu),
		Places: places,
	}
}
