package handler

import (
	"net/http"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/groups", func() trama.WebHandler {
		return new(groups)
	})
}

type groups struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
}

func (h *groups) Get(response trama.Response, r *http.Request) {
	groupDAO := dao.NewGroup(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	groups, err := groupDAO.FindAll()

	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate("groups.html",
		data.NewGroups(h.Session().User.Name, data.MenuGroups, groups))
}

func (h *groups) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(nil)

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "groups")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *groups) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
	)
}
