package payment

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

func (s Service) Save(ip net.IP, agent int, p *Payment) error {
	dao := newDAO(s.sqler, ip, agent)
	return dao.save(p)
}

func (s Service) FindById(id int) (Payment, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindByClient(clientId int) ([]Payment, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findByClient(clientId)
}
