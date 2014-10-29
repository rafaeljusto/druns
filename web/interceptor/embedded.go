package interceptor

import (
	"reflect"

	"github.com/rafaeljusto/druns/core/protocol"
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

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type JSONCompliant struct {
	request  reflect.Value
	response reflect.Value
	message  protocol.Translator
}

func (j *JSONCompliant) RequestValue() reflect.Value {
	return j.request
}

func (j *JSONCompliant) SetRequestValue(r reflect.Value) {
	j.request = r
}

func (j *JSONCompliant) ResponseValue() reflect.Value {
	return j.response
}

func (j *JSONCompliant) SetResponseValue(r reflect.Value) {
	j.response = r
}

func (j *JSONCompliant) Message() protocol.Translator {
	return j.message
}

func (j *JSONCompliant) SetMessage(m protocol.Translator) {
	j.message = m
}
