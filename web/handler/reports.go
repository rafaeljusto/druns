package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/reports"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/reports", func() trama.WebHandler {
		return new(reportsHandler)
	})
}

type reportsHandler struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
}

func (h *reportsHandler) Get(response trama.Response, r *http.Request) {
	incomings, err := reports.NewService(h.Tx()).IncomingPerGroup(time.Now(), config.DrunsConfig.ClassValue)

	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate("reports.html",
		data.NewReports(h.Session().User.Name, data.MenuReports, incomings))
}

func (h *reportsHandler) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(template.FuncMap{
		"month": func(time time.Time) string {
			return fmt.Sprintf("%02d/%d", time.Month(), time.Year())
		},
	})

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "reports")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *reportsHandler) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
	)
}
