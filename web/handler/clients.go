package handler

import (
	"net/http"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/clients", func() trama.WebHandler {
		return new(clients)
	})
}

type clients struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
}

func (h *clients) Get(response trama.Response, r *http.Request) {
	clientDAO := dao.NewClient(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	clients, err := clientDAO.FindAll()

	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate("clients.html",
		data.NewClients(h.Session().User.Name, data.MenuClients, clients))
}

func (h *clients) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(nil)

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "clients")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *clients) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
	)
}
