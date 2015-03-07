package dblog

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

func (s Service) Create(agent int, ipAddress net.IP, operation Operation) (DBLog, error) {
	dbLog := NewDBLog(agent, ipAddress, operation)
	dao := newDAO(s.sqler)
	return dbLog, dao.save(&dbLog)
}
