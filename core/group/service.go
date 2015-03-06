package group

import (
	"net"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) Save(sqler db.SQLer, ip net.IP, agent int, g *Group) error {
	dao := newDAO(sqler, ip, agent)
	return dao.save(g)
}

func (s Service) FindById(sqler db.SQLer, id int) (Group, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindAll(sqler db.SQLer) (Groups, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.findAll()
}
