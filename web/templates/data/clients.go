package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/client"
)

type Clients struct {
	Logged
	Clients []client.Client
}

func NewClients(username core.Name, menu Menu, clients []client.Client) Clients {
	return Clients{
		Logged:  NewLogged(username, menu),
		Clients: clients,
	}
}
