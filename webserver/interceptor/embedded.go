package interceptor

import (
	"gopkg.in/mgo.v2"
)

type DatabaseCompliant struct {
	db      *mgo.Database
	session *mgo.Session
}

func (d *DatabaseCompliant) SetDBSession(s *mgo.Session) {
	d.session = s
}

func (d *DatabaseCompliant) DBSession() *mgo.Session {
	return d.session
}

func (d *DatabaseCompliant) SetDB(db *mgo.Database) {
	d.db = db
}

func (d *DatabaseCompliant) DB() *mgo.Database {
	return d.db
}
