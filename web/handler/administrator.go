package handler

import (
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/model"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
	"github.com/rafaeljusto/druns/web/tr"
)

func init() {
	Mux.RegisterPage("/administrator", func() trama.WebHandler {
		return new(administrator)
	})
}

type administrator struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
}

func (h *administrator) Get(response trama.Response, r *http.Request) {
	if len(r.FormValue("id")) == 0 {
		response.ExecuteTemplate("administrator.html",
			data.NewAdministrator(h.Session().User.Name, data.MenuAdministrators, model.User{}))
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	userDAO := dao.NewUser(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	user, err := userDAO.FindById(id)

	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate("administrator.html",
		data.NewAdministrator(h.Session().User.Name, data.MenuAdministrators, user))
}

func (h *administrator) Post(response trama.Response, r *http.Request) {
	user := model.User{}

	if len(r.FormValue("id")) > 0 {
		var err error
		user.Id, err = strconv.Atoi(r.FormValue("id"))
		if err != nil {
			h.Logger().Error(err)
			response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
			return
		}
	}

	user.Name = r.FormValue("name")
	user.Name = strings.TrimSpace(user.Name)
	user.Name = strings.Title(user.Name)

	user.Email = r.FormValue("email")
	user.Email = strings.TrimSpace(user.Email)
	user.Email = strings.ToLower(user.Email)

	if _, err := mail.ParseAddress(user.Email); err != nil {
		data := data.NewAdministrator(h.Session().User.Name, data.MenuAdministrators, user)
		data.FieldMessage["email"] = h.Msg(tr.CodeInvalidEmail)
		response.ExecuteTemplate("administrator.html", data)
		return
	}

	userDAO := dao.NewUser(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	if err := userDAO.Save(&user); err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("administrators"), http.StatusFound)
	return
}

func (h *administrator) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(nil)

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "administrator")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *administrator) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
	)
}
