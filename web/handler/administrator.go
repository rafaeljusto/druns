package handler

import (
	"net/http"
	"strconv"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/model"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
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
	interceptor.POSTCompliant

	User model.User `request:"post"`
}

func (h administrator) Response() (string, data.Former) {
	data := data.NewAdministrator(h.Session().User.Name, data.MenuAdministrators)
	data.User = h.User
	return "administrator.html", &data
}

func (h *administrator) Get(response trama.Response, r *http.Request) {
	if len(r.FormValue("id")) == 0 {
		response.ExecuteTemplate(h.Response())
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		h.Logger().Error(core.NewError(err))
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	userDAO := dao.NewUser(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	h.User, err = userDAO.FindById(id)

	if err != nil {
		// TODO: Check ErrNotFound. Redirect to the list page with an automatic error message (like login)

		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate(h.Response())
}

func (h *administrator) Post(response trama.Response, r *http.Request) {
	userDAO := dao.NewUser(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	if err := userDAO.Save(&h.User); err != nil {
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
		interceptor.NewPOST(h),
	)
}
