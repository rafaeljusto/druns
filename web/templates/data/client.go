package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/client"
)

type Client struct {
	Logged
	Form
	Client client.Client
}

func NewClient(username core.Name, menu Menu) Client {
	return Client{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
