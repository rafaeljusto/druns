package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/client"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/enrollment", func() trama.WebHandler {
		return new(enrollmentHandler)
	})
}

type enrollmentHandler struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
	interceptor.POSTCompliant

	Enrollment enrollment.Enrollment `request:"post"`
}

func (h enrollmentHandler) Response(r *http.Request) (string, data.Former) {
	data := data.NewEnrollment(h.Session().User.Name, data.MenuGroups)
	data.Enrollment = h.Enrollment
	data.Clients, _ = client.NewService().FindAll(h.Tx())
	data.Groups, _ = group.NewService().FindAll(h.Tx())
	data.Back = r.FormValue("back")
	return "enrollment.html", &data
}

func (h *enrollmentHandler) Get(response trama.Response, r *http.Request) {
	if len(r.FormValue("id")) == 0 {
		response.ExecuteTemplate(h.Response(r))
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		h.Logger().Error(core.NewError(err))
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	if h.Enrollment, err = enrollment.NewService().FindById(h.Tx(), id); err != nil {
		// TODO: Check ErrNotFound. Redirect to the list page with an automatic error message (like login)
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate(h.Response(r))
}

func (h *enrollmentHandler) Post(response trama.Response, r *http.Request) {
	err := enrollment.NewService().Save(h.Tx(), h.RemoteAddress(), h.Session().User.Id, &h.Enrollment)
	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	back := r.FormValue("back")
	if len(back) == 0 {
		back = fmt.Sprintf("%s?id=%d", config.DrunsConfig.URLs.GetHTTPS("group"), h.Enrollment.Group.Id)
	}

	response.Redirect(back, http.StatusFound)
	return
}

func (h *enrollmentHandler) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(template.FuncMap{
		"teq": func(value1 enrollment.Type, value2 string) bool {
			return value1.String() == value2
		},
	})

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "enrollment")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *enrollmentHandler) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
		interceptor.NewPOST(h),
	)
}
