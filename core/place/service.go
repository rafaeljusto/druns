package place

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

func (s Service) Save(ip net.IP, agent int, p *Place) error {
	dao := newDAO(s.sqler, ip, agent)
	return dao.save(p)
}

func (s Service) FindById(id int) (Place, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindAll() ([]Place, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findAll()
}
