package user

import (
	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) FindById(sqler db.SQLer, id int) (User, error) {
	dao := newDAO(sqler, nil, 0)
	return dao.FindById(id)
}
