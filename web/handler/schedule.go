package handler

import (
	"net/http"
	"time"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/class"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/interceptor"
	"github.com/rafaeljusto/druns/web/templates/data"
)

func init() {
	Mux.RegisterPage("/schedule", func() trama.WebHandler {
		return new(schedule)
	})
}

type schedule struct {
	trama.DefaultWebHandler
	interceptor.DatabaseCompliant
	interceptor.RemoteAddressCompliant
	interceptor.LanguageCompliant
	interceptor.HTTPTransactionCompliant
	interceptor.SessionCompliant
}

func (h *schedule) Get(response trama.Response, r *http.Request) {
	var begin time.Time

	if r.FormValue("begin") != "" {
		var err error
		if begin, err = time.Parse("2006-01-02", r.FormValue("begin")); err != nil {
			h.Logger().Error(errors.New(err))
			response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
			return
		}

	} else {
		begin = time.Now()
		for begin.Weekday() != time.Sunday {
			begin = begin.Add(time.Duration(-24) * time.Hour)
		}
	}

	end := begin.Add(time.Duration(7*24) * time.Hour)
	classes, err := class.NewClassService(h.Tx()).FindBetweenDates(begin, end)

	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate("schedule.html",
		data.NewSchedule(h.Session().User.Name, data.MenuSchedule, classes))
}

func (h *schedule) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(nil)

	for _, language := range config.DrunsConfig.Languages {
		templates := config.DrunsConfig.HTMLTemplates(language, "schedule")

		groupSet.Insert(trama.TemplateGroup{
			Name:  language,
			Files: templates,
		})
	}

	return groupSet
}

func (h *schedule) Interceptors() trama.WebInterceptorChain {
	return trama.NewWebInterceptorChain(
		interceptor.NewRemoteAddressWeb(h),
		interceptor.NewAcceptLanguageWeb(h),
		interceptor.NewHTTPTransactionWeb(h),
		interceptor.NewDatabaseWeb(h),
		interceptor.NewSessionWeb(h),
	)
}
