package interceptor

import (
	"net"
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/session"
)

type auther interface {
	Tx() db.Transaction
	RemoteAddress() net.IP
	Logger() *log.Logger
}

////////////////////////////////////////////////////////////
/////////////////////// AJAX ///////////////////////////////
////////////////////////////////////////////////////////////

type AuthAJAX struct {
	trama.NopAJAXInterceptor
	handler auther
}

func NewAuthAJAX(h auther) *AuthAJAX {
	return &AuthAJAX{handler: h}
}

func (i AuthAJAX) Before(w http.ResponseWriter, r *http.Request) {
	ok, err := auth(r, i.handler.Tx(), i.handler.RemoteAddress())
	if err == nil && ok {
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

type AuthWeb struct {
	trama.NopWebInterceptor
	handler auther
}

func NewAuthWeb(h auther) *AuthWeb {
	return &AuthWeb{handler: h}
}

func (i AuthWeb) Before(response trama.Response, r *http.Request) {
	ok, err := auth(r, i.handler.Tx(), i.handler.RemoteAddress())
	if err == nil && ok {
		return
	}

	if err != nil {
		i.handler.Logger().Error(err)
	}

	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("login"), http.StatusFound)
}

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

func auth(r *http.Request, tx db.Transaction, remoteAddress net.IP) (bool, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false, core.NewError(err)
	}

	return session.CheckSession(tx, cookie, remoteAddress)
}
