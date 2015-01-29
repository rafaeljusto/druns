package interceptor

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/db"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/session"
)

////////////////////////////////////////////////////////////
/////////////////////// AJAX ///////////////////////////////
////////////////////////////////////////////////////////////

type AuthAJAX struct {
	trama.NopAJAXInterceptor
	handler sqler
}

func NewAuthAJAX(h sqler) *AuthAJAX {
	return &AuthAJAX{handler: h}
}

func (i AuthAJAX) Before(response trama.Response, r *http.Request) {
	ok, err := auth(r, i.handler.Tx())
	if err == nil && ok {
		return
	}

	if err != nil {
		i.handler.Logger().Error(err)
	}

	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("login"), http.StatusFound)
}

////////////////////////////////////////////////////////////
/////////////////////// WEB ////////////////////////////////
////////////////////////////////////////////////////////////

type AuthWeb struct {
	trama.NopWebInterceptor
	handler sqler
}

func NewAuthWeb(h sqler) *AuthWeb {
	return &AuthWeb{handler: h}
}

func (i AuthWeb) Before(response trama.Response, r *http.Request) {
	ok, err := auth(r, i.handler.Tx())
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

func auth(r *http.Request, tx db.Transaction) (bool, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false, err
	}

	return session.CheckSession(tx, cookie)
}
