package class

import (
	"net"
	"time"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
	sqler db.SQLer
}

func NewService(sqler db.SQLer) Service {
	return Service{sqler}
}

func (s Service) Save(ip net.IP, agent int, c *Class) error {
	dao := newDAO(s.sqler, ip, agent)
	return dao.save(c)
}

func (s Service) FindById(id int) (Class, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindAll() ([]Class, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findAll()
}

func (s Service) FindByGroupIdBetweenDates(groupId int, begin, end time.Time) ([]Class, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findByGroupIdBetweenDates(groupId, begin, end)
}
