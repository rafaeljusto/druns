package enrollment

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

func (s Service) Save(ip net.IP, agent int, e *Enrollment) error {
	dao := newDAO(s.sqler, ip, agent)
	return dao.save(e)
}

func (s Service) FindById(id int) (Enrollment, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindByClient(clientId int) ([]Enrollment, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findByClient(clientId)
}

func (s Service) FindByGroup(groupId int) ([]Enrollment, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findByGroup(groupId)
}
