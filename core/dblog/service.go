package dblog

import (
	"net"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
}

func NewService() Service {
	return Service{}
}

func (s Service) Create(sqler db.SQLer, agent int, ipAddress net.IP, operation Operation) (DBLog, error) {
	dbLog := NewDBLog(agent, ipAddress, operation)
	dao := newDAO(sqler)
	return dbLog, dao.save(&dbLog)
}
