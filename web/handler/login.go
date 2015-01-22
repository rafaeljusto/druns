package handler

import (
	"net/http"
	"net/mail"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/web/interceptor"
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
	response.ExecuteTemplate("login.html", nil)
}

func (h *login) Post(response trama.Response, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	address, err := mail.ParseAddress(email)
	if err != nil {
		// TODO
	}

	userDAO := dao.NewUser(h.Tx(), h.RemoteAddress(), "")
	if ok, err := userDAO.VerifyPassword(*address, password); !ok || err != nil {
		// TODO
	} else {
		// TODO
	}
}

func (h *login) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
	)
}
