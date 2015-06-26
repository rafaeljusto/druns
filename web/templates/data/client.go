package data

import (
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/payment"
	"github.com/rafaeljusto/druns/core/types"
)

type Client struct {
	Logged
	Form
	Client      client.Client
	Enrollments []enrollment.Enrollment
	Payments    []payment.Payment
}

func NewClient(username types.Name, menu Menu) Client {
	return Client{
		Logged: NewLogged(username, menu),
		Form:   NewForm(),
	}
}
