package place

import (
	"net"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) Save(sqler db.SQLer, ip net.IP, agent int, p *Place) error {
	dao := newDAO(sqler, ip, agent)
	return dao.save(p)
}

func (s Service) FindById(sqler db.SQLer, id int) (Place, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindAll(sqler db.SQLer) (Places, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.findAll()
}
