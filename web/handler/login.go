package handler

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/tr"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/session"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/", func() trama.WebHandler {
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
		_, err := session.LoadAndCheckSession(h.Tx(), cookie, h.RemoteAddress())
		if err == nil {
			response.Redirect(config.DrunsConfig.URLs.GetHTTPS("schedule"), http.StatusFound)
			return
		}
	}

	data := data.NewLogin("")
	if message := r.FormValue("m"); len(message) > 0 {
		data.Message = h.Msg(tr.Code(message))
	}

	response.ExecuteTemplate("login.html", data)
}

func (h *login) Post(response trama.Response, r *http.Request) {
	email := r.FormValue("email")
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)

	password := r.FormValue("password")

	address, err := mail.ParseAddress(email)
	if err != nil {
		data := data.NewLogin(email)
		data.FieldMessage["email"] = h.Msg(tr.CodeInvalidEmail)
		response.ExecuteTemplate("login.html", data)
		return
	}

	userDAO := dao.NewUser(h.Tx(), h.RemoteAddress(), 0)
	if ok, err := userDAO.VerifyPassword(*address, password); !ok || err != nil {
		if err != nil {
			h.Logger().Error(err)
		}
		data := data.NewLogin(email)
		data.Message = h.Msg(tr.CodeAuthenticationError)
		response.ExecuteTemplate("login.html", data)
		return
	}

	cookie, err := session.NewSession(h.Tx(), email, h.RemoteAddress())
	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.SetCookie(cookie)
	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("schedule"), http.StatusFound)
}

func (h *login) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(nil)

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "login")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *login) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
	)
}
