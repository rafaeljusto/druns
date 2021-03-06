package data

import (
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/types"
)

type Clients struct {
	Logged
	Clients []client.Client
}

func NewClients(username types.Name, menu Menu, clients []client.Client) Clients {
	return Clients{
		Logged:  NewLogged(username, menu),
		Clients: clients,
	}
}
