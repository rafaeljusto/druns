package data

import (
	"github.com/rafaeljusto/druns/core/place"
	"github.com/rafaeljusto/druns/core/types"
)

type Places struct {
	Logged
	Places []place.Place
}

func NewPlaces(username types.Name, menu Menu, places []place.Place) Places {
	return Places{
		Logged: NewLogged(username, menu),
		Places: places,
	}
}
