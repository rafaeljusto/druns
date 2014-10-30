package interceptor

import (
	"net"
	"reflect"

	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/core/protocol"
	"github.com/rafaeljusto/druns/web/tr"
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

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type LanguageCompliant struct {
	language string
	messages tr.MessageHolder
}

func (h *LanguageCompliant) SetLanguage(language string) {
	h.language = language
}

func (h *LanguageCompliant) Language() string {
	return h.language
}

func (h *LanguageCompliant) SetMessages(messages tr.MessageHolder) {
	h.messages = messages
}

func (h *LanguageCompliant) Msg(code tr.Code, args ...interface{}) string {
	return h.messages.Get(code, args...)
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type RemoteAddressCompliant struct {
	remoteAddress net.IP
}

func (r *RemoteAddressCompliant) SetRemoteAddress(a net.IP) {
	r.remoteAddress = a
}

func (r *RemoteAddressCompliant) RemoteAddress() net.IP {
	return r.remoteAddress
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type LogCompliant struct {
	logger log.Logger
}

func (l *LogCompliant) SetLogger(logger log.Logger) {
	l.logger = logger
}

func (l *LogCompliant) Logger() *log.Logger {
	return &l.logger
}
