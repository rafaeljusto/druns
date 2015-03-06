package session

import (
	"net"

	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/user"
)

type Service struct {
	sqler db.SQLer
}

func NewService(sqler db.SQLer) Service {
	return Service{sqler}
}

func (s Service) Create(u user.User, ip net.IP) (Session, error) {
	session := NewSession(u, ip)
	dao := newDAO(s.sqler)
	return session, dao.save(&session)
}

func (s Service) Save(session *Session) error {
	dao := newDAO(s.sqler)
	return dao.save(session)
}

func (s Service) FindById(id int) (Session, error) {
	dao := newDAO(s.sqler)
	return dao.findById(id)
}
