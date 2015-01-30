package interceptor

import (
	"net"
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/core/model"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/session"
	"github.com/rafaeljusto/druns/web/tr"
)

type sessioner interface {
	Tx() db.Transaction
	RemoteAddress() net.IP
	Logger() *log.Logger
	SetSession(session model.Session)
	Session() model.Session
}

////////////////////////////////////////////////////////////
/////////////////////// AJAX ///////////////////////////////
////////////////////////////////////////////////////////////

type SessionAJAX struct {
	trama.NopAJAXInterceptor
	handler sessioner
}

func NewSessionAJAX(h sessioner) *SessionAJAX {
	return &SessionAJAX{handler: h}
}

func (i SessionAJAX) Before(w http.ResponseWriter, r *http.Request) {
	session, err := auth(r, i.handler.Tx(), i.handler.RemoteAddress())
	if err == nil {
		i.handler.SetSession(session)
		return
	}

	if err != nil {
		i.handler.Logger().Error(err)
	}

	w.WriteHeader(http.StatusUnauthorized)
}

////////////////////////////////////////////////////////////
/////////////////////// WEB ////////////////////////////////
////////////////////////////////////////////////////////////

type SessionWeb struct {
	trama.NopWebInterceptor
	handler sessioner
}

func NewSessionWeb(h sessioner) *SessionWeb {
	return &SessionWeb{handler: h}
}

func (i SessionWeb) Before(response trama.Response, r *http.Request) {
	session, err := auth(r, i.handler.Tx(), i.handler.RemoteAddress())
	if err == nil {
		i.handler.SetSession(session)
		return
	}

	if err != nil {
		i.handler.Logger().Error(err)
	}

	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("login", "m="+string(tr.CodeSessionExpired)),
		http.StatusFound)
}

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

func auth(r *http.Request, tx db.Transaction, remoteAddress net.IP) (model.Session, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return model.Session{}, core.NewError(err)
	}

	return session.LoadAndCheckSession(tx, cookie, remoteAddress)
}
