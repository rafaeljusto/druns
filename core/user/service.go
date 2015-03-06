package user

import (
	"net"
	"net/mail"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
	sqler db.SQLer
}

func NewService(sqler db.SQLer) Service {
	return Service{sqler}
}

func (s Service) Save(ip net.IP, agent int, u *User) error {
	dao := newDAO(s.sqler, ip, agent)
	return dao.save(u)
}

func (s Service) FindById(id int) (User, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findById(id)
}

func (s Service) FindByEmail(email string) (User, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findByEmail(email)
}

func (s Service) FindAll() ([]User, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.findAll()
}

func (s Service) VerifyPassword(email mail.Address, password string) (bool, error) {
	dao := newDAO(s.sqler, nil, 0)
	return dao.verifyPassword(email, password)
}
