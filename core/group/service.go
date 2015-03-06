package group

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

func (s Service) Save(ip net.IP, agent int, g *Group) error {
	dao := newDAO(s.sqler, ip, agent)
	return dao.save(g)
}

func (s Service) FindById(id int) (Group, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindAll() (Groups, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findAll()
}
