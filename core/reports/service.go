package reports

import (
	"time"

	"github.com/rafaeljusto/druns/core/db"
)

type Service struct {
	sqler db.SQLer
}

func NewService(sqler db.SQLer) Service {
	return Service{sqler}
}

func (s Service) IncomingPerGroup(month time.Time, classValue float64) ([]Incoming, error) {
	dao := newDAO(s.sqler)
	return dao.incomingPerGroup(month, classValue)
}
