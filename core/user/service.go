package user

import (
	"net"
	"net/mail"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
	sqler db.SQLer
	ip    net.IP
	agent int
}

func NewService() Service {
	return Service{}
}

func (s Service) Save(sqler db.SQLer, ip net.IP, agent int, u *User) error {
	dao := newDAO(sqler, ip, agent)
	return dao.Save(u)
}

func (s Service) FindById(sqler db.SQLer, id int) (User, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.FindById(id)
}

func (s Service) FindByEmail(sqler db.SQLer, email string) (User, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.FindByEmail(email)
}

func (s Service) FindAll(sqler db.SQLer) ([]User, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.FindAll()
}

func (s Service) VerifyPassword(sqler db.SQLer, email mail.Address, password string) (bool, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.VerifyPassword(email, password)
}
