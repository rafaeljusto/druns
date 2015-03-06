package client

import (
	"net"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
	sqler db.SQLer
}

func NewService(sqler db.SQLer) Service {
	return Service{sqler}
}

func (s Service) Save(ip net.IP, agent int, c *Client) error {
	dao := newDAO(s.sqler, ip, agent)
	return dao.save(c)
}

func (s Service) FindById(id int) (Client, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindAll() (Clients, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findAll()
}
