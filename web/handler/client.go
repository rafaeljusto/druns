package handler

import (
	"net/http"
	"strconv"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/client", func() trama.WebHandler {
		return new(clientHandler)
	})
}

type clientHandler struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
	interceptor.POSTCompliant

	Client client.Client `request:"post"`
}

func (h clientHandler) Response(r *http.Request) (string, data.Former) {
	data := data.NewClient(h.Session().User.Name, data.MenuClients)
	data.Client = h.Client

	if h.Client.Id > 0 {
		var err error
		data.Enrollments, err = enrollment.NewService().FindByClient(h.Tx(), h.Client.Id)
		if err != nil {
			h.Logger().Error(errors.New(err))
		}
	}

	return "client.html", &data
}

func (h *clientHandler) Get(response trama.Response, r *http.Request) {
	if len(r.FormValue("id")) == 0 {
		response.ExecuteTemplate(h.Response(r))
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		h.Logger().Error(errors.New(err))
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	if h.Client, err = client.NewService().FindById(h.Tx(), id); err != nil {
		// TODO: Check ErrNotFound. Redirect to the list page with an automatic error message (like login)
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate(h.Response(r))
}

func (h *clientHandler) Post(response trama.Response, r *http.Request) {
	err := client.NewService().Save(h.Tx(), h.RemoteAddress(), h.Session().User.Id, &h.Client)
	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("clients"), http.StatusFound)
	return
}

func (h *clientHandler) Templates() trama.TemplateGroupSet {
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

func (h *clientHandler) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
		interceptor.NewPOST(h),
	)
}
