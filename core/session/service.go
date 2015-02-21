package session

import (
	"net"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/user"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) Create(sqler db.SQLer, u user.User, ip net.IP) (Session, error) {
	session := NewSession(u, ip)
	dao := newDAO(sqler)
	return session, dao.Save(&session)
}

func (s Service) Save(sqler db.SQLer, session *Session) error {
	dao := newDAO(sqler)
	return dao.Save(session)
}

func (s Service) FindById(sqler db.SQLer, id int) (Session, error) {
	dao := newDAO(sqler)
	return dao.FindById(id)
}
