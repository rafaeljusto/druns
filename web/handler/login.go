package handler

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/session"
	"github.com/rafaeljusto/druns/web/templates/data"
	"github.com/rafaeljusto/druns/web/tr"
)

func init() {
	Mux.RegisterPage("/login", func() trama.WebHandler {
		return new(login)
	})
}

type login struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
}

func (h *login) Get(response trama.Response, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		ok, err := session.CheckSession(h.Tx(), cookie, h.RemoteAddress())
		if err == nil && ok {
			response.Redirect(config.DrunsConfig.URLs.GetHTTPS("home"), http.StatusFound)
			return
		}
	}

	response.ExecuteTemplate("login.html", data.NewLogin("", ""))
}

func (h *login) Post(response trama.Response, r *http.Request) {
	email := r.FormValue("email")
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	password := r.FormValue("password")

	address, err := mail.ParseAddress(email)
	if err != nil {
		response.ExecuteTemplate("login.html", data.NewLogin(email, h.Msg(tr.CodeInvalidEmail)))
		return
	}

	userDAO := dao.NewUser(h.Tx(), h.RemoteAddress(), "")
	if ok, err := userDAO.VerifyPassword(*address, password); !ok || err != nil {
		if err != nil {
			h.Logger().Error(err)
		}
		response.ExecuteTemplate("login.html", data.NewLogin(email, h.Msg(tr.CodeAuthenticationError)))
		return
	}

	cookie, err := session.NewSession(h.Tx(), email, h.RemoteAddress())
	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.SetCookie(cookie)
	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("home"), http.StatusFound)
}

func (h *login) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
	)
}
