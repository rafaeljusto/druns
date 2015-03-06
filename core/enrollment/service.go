package enrollment

import (
	"net"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) Save(sqler db.SQLer, ip net.IP, agent int, e *Enrollment) error {
	dao := newDAO(sqler, ip, agent)
	return dao.save(e)
}

func (s Service) FindById(sqler db.SQLer, id int) (Enrollment, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindByClient(sqler db.SQLer, clientId int) (Enrollments, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.findByClient(clientId)
}

func (s Service) FindByGroup(sqler db.SQLer, groupId int) (Enrollments, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.findByGroup(groupId)
}
