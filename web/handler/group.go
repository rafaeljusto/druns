package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/enrollment"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/group"
	"github.com/rafaeljusto/druns/core/place"
	"github.com/rafaeljusto/druns/core/types"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/group", func() trama.WebHandler {
		return new(groupHandler)
	})
}

type groupHandler struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
	interceptor.POSTCompliant

	Group group.Group `request:"post"`
}

func (h groupHandler) Response(r *http.Request) (string, data.Former) {
	data := data.NewGroup(h.Session().User.Name, data.MenuGroups)
	data.Group = h.Group

	var err error
	data.Places, err = place.NewService().FindAll(h.Tx())
	if err != nil {
		h.Logger().Error(errors.New(err))
	}

	if h.Group.Id > 0 {
		data.Enrollments, err = enrollment.NewService().FindByGroup(h.Tx(), h.Group.Id)
		if err != nil {
			h.Logger().Error(errors.New(err))
		}
	}

	return "group.html", &data
}

func (h *groupHandler) Get(response trama.Response, r *http.Request) {
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

	if h.Group, err = group.NewService().FindById(h.Tx(), id); err != nil {
		// TODO: Check ErrNotFound. Redirect to the list page with an automatic error message (like login)
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate(h.Response(r))
}

func (h *groupHandler) Post(response trama.Response, r *http.Request) {
	err := group.NewService().Save(h.Tx(), h.RemoteAddress(), h.Session().User.Id, &h.Group)
	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	// TODO: Maybe we should keep the user on this page to allow him to add
	// enrollments. If so, we need to add a success message to make it clear
	// that the object was created or updated
	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("groups"), http.StatusFound)
	return
}

func (h *groupHandler) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(template.FuncMap{
		"weq": func(value1 types.Weekday, value2 string) bool {
			return value1.String() == value2
		},
		"geq": func(value1 group.Type, value2 string) bool {
			return value1.String() == value2
		},
	})

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "group")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *groupHandler) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
		interceptor.NewPOST(h),
	)
}
