package data

import "github.com/rafaeljusto/druns/core/model"

type Client struct {
	Logged
	Form
	Client model.Client
}

func NewClient(username model.Name, menu Menu) Client {
	return Client{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
