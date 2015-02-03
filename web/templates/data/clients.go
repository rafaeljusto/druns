package data

import "github.com/rafaeljusto/druns/core/model"

type Clients struct {
	Logged
	Clients []model.Client
}

func NewClients(username string, menu Menu, clients []model.Client) Clients {
	return Clients{
		Logged:  NewLogged(username, menu),
		Clients: clients,
	}
}
