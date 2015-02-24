package data

import (
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/enrollment"
)

type Client struct {
	Logged
	Form
	Client      client.Client
	Enrollments []enrollment.Enrollment
}

func NewClient(username core.Name, menu Menu) Client {
	return Client{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
