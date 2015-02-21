package client

import (
	"net"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) Save(sqler db.SQLer, ip net.IP, agent int, c *Client) error {
	dao := newDAO(sqler, ip, agent)
	return dao.Save(c)
}

func (s Service) FindById(sqler db.SQLer, id int) (Client, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.FindById(id)
}

func (s Service) FindAll(sqler db.SQLer) (Clients, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.FindAll()
}
