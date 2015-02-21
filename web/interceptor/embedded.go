package interceptor

import (
	"net"
	"reflect"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/core/protocol"
	"github.com/rafaeljusto/druns/core/session"
	"github.com/rafaeljusto/druns/web/tr"
)

type DatabaseCompliant struct {
	tx db.Transaction
}

func (d *DatabaseCompliant) Tx() db.Transaction {
	return d.tx
}

func (d *DatabaseCompliant) SetTx(tx db.Transaction) {
	d.tx = tx
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

type POSTCompliant struct {
	request reflect.Value
}

func (p *POSTCompliant) RequestValue() reflect.Value {
	return p.request
}

func (p *POSTCompliant) SetRequestValue(r reflect.Value) {
	p.request = r
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

func (h *LanguageCompliant) Msg(code core.ValidationErrorCode, args ...interface{}) string {
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

type HTTPTransactionCompliant struct {
	logger log.Logger
	id     string
}

func (h *HTTPTransactionCompliant) SetLogger(logger log.Logger) {
	h.logger = logger
}

func (h *HTTPTransactionCompliant) Logger() *log.Logger {
	return &h.logger
}

func (h *HTTPTransactionCompliant) SetHTTPId(id string) {
	h.id = id
}

func (h *HTTPTransactionCompliant) HTTPId() string {
	return h.id
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type SessionCompliant struct {
	session session.Session
}

func (h *SessionCompliant) SetSession(session session.Session) {
	h.session = session
}

func (h *SessionCompliant) Session() session.Session {
	return h.session
}
