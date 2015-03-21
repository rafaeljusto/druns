package handler

import (
	"html/template"
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

	end := begin.Add(time.Duration(6*24) * time.Hour)
	classes, err := class.NewClassService(h.Tx()).FindBetweenDates(begin, end)

	if err != nil {
		h.Logger().Error(err)
		response.ExecuteTemplate("500.html", data.NewInternalServerError(h.HTTPId()))
		return
	}

	response.ExecuteTemplate("schedule.html",
		data.NewSchedule(h.Session().User.Name, data.MenuSchedule, begin, end, classes))
}

func (h *schedule) Templates() trama.TemplateGroupSet {
	groupSet := trama.NewTemplateGroupSet(template.FuncMap{
		"getWeekdays": func(begin, end time.Time) []time.Time {
			var days []time.Time
			for begin.Day() <= end.Day() &&
				begin.Month() <= end.Month() &&
				begin.Year() <= end.Year() {

				days = append(days, begin)
				begin = begin.Add(24 * time.Hour)
			}
			return days
		},
		"weekday": func(date time.Time) string {
			switch date.Weekday() {
			case time.Sunday:
				return h.Msg(errors.ValidationCodeSunday)
			case time.Monday:
				return h.Msg(errors.ValidationCodeMonday)
			case time.Tuesday:
				return h.Msg(errors.ValidationCodeTuesday)
			case time.Wednesday:
				return h.Msg(errors.ValidationCodeWednesday)
			case time.Thursday:
				return h.Msg(errors.ValidationCodeThursday)
			case time.Friday:
				return h.Msg(errors.ValidationCodeFriday)
			case time.Saturday:
				return h.Msg(errors.ValidationCodeSaturday)
			}

			return ""
		},
		"printHour": func(hour time.Time) string {
			return hour.Format("15:04")
		},
		"getWorkingHours": func() []time.Time {
			var workingHours []time.Time
			for i := 6; i <= 23; i++ {
				workingHours = append(workingHours, time.Date(0, 0, 0, i, 0, 0, 0, time.UTC))
			}
			return workingHours
		},
		"filterClasses": func(classes []class.Class, date time.Time, hour time.Time) []class.Class {
			begin := time.Date(date.Year(), date.Month(), date.Day(), hour.Hour(), 0, 0, 0, time.UTC)
			end := begin.Add(1 * time.Hour)

			var filtered []class.Class
			for _, c := range classes {
				beginInDate := (c.BeginAt.After(begin) || c.BeginAt.Equal(begin)) &&
					(c.BeginAt.Before(end) || c.BeginAt.Equal(end))
				endInDate := (c.EndAt.After(begin) || c.EndAt.Equal(begin)) &&
					(c.EndAt.Before(end) || c.EndAt.Equal(end))

				if beginInDate || endInDate {
					filtered = append(filtered, c)
				}
			}
			return filtered
		},
	})

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
