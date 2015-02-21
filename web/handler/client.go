package handler

import (
	"net/http"
	"strconv"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/dao"
	"github.com/rafaeljusto/druns/core/model"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/client", func() trama.WebHandler {
		return new(client)
	})
}

type client struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
	interceptor.POSTCompliant

	Client model.Client `request:"post"`
}

func (h client) Response() (string, data.Former) {
	data := data.NewClient(h.Session().User.Name, data.MenuClients)
	data.Client = h.Client
	return "client.html", &data
}

func (h *client) Get(response trama.Response, r *http.Request) {
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

	clientDAO := dao.NewClient(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	h.Client, err = clientDAO.FindById(id)

	if err != nil {
		// TODO: Check ErrNotFound. Redirect to the list page with an automatic error message (like login)

		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate(h.Response())
}

func (h *client) Post(response trama.Response, r *http.Request) {
	clientDAO := dao.NewClient(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	if err := clientDAO.Save(&h.Client); err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("clients"), http.StatusFound)
	return
}

func (h *client) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(nil)

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "client")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *client) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
		interceptor.NewPOST(h),
	)
}
