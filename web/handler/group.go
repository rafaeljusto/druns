package handler

import (
	"html/template"
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
	Mux.RegisterPage("/group", func() trama.WebHandler {
		return new(group)
	})
}

type group struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
	interceptor.POSTCompliant

	Group model.Group `request:"post"`
}

func (h group) Response() (string, data.Former) {
	placeDAO := dao.NewPlace(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	places, _ := placeDAO.FindAll()

	data := data.NewGroup(h.Session().User.Name, data.MenuGroups)
	data.Group = h.Group
	data.Places = places

	return "group.html", &data
}

func (h *group) Get(response trama.Response, r *http.Request) {
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

	groupDAO := dao.NewGroup(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	h.Group, err = groupDAO.FindById(id)

	if err != nil {
		// TODO: Check ErrNotFound. Redirect to the list page with an automatic error message (like login)

		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate(h.Response())
}

func (h *group) Post(response trama.Response, r *http.Request) {
	groupDAO := dao.NewGroup(h.Tx(), h.RemoteAddress(), h.Session().User.Id)
	if err := groupDAO.Save(&h.Group); err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.Redirect(config.DrunsConfig.URLs.GetHTTPS("groups"), http.StatusFound)
	return
}

func (h *group) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(template.FuncMap{
		"weq": func(value1 model.Weekday, value2 string) bool {
			return value1.String() == value2
		},
		"geq": func(value1 model.GroupType, value2 string) bool {
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

func (h *group) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
		interceptor.NewPOST(h),
	)
}
