package interceptor

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/rafaeljusto/schema"
	"github.com/rafaeljusto/druns/core/errors"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/web/templates/data"
)

type poster interface {
	RequestValue() reflect.Value
	SetRequestValue(reflect.Value)
	Response(r *http.Request) (string, data.Former)
	Msg(code errors.ValidationCode, args ...interface{}) string
	Logger() *log.Logger
	HTTPId() string
}

type POST struct {
	trama.NopWebInterceptor
	handler poster
}

func NewPOST(h poster) *POST {
	return &POST{handler: h}
}

func (i *POST) Before(response trama.Response, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	i.parse()

	request := i.handler.RequestValue()
	if !request.IsValid() {
		return
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := r.ParseForm(); err != nil {
		i.handler.Logger().Error(errors.New(err))
		response.ExecuteTemplate("500.html", data.NewInternalServerError(i.handler.HTTPId()))
		return
	}

	if request.CanAddr() {
		request = request.Addr()
	}

	err := decoder.Decode(request.Interface(), r.Form)
	if err == nil {
		return
	}

	if conversionErr, ok := err.(schema.ConversionError); ok {
		template, former := i.handler.Response(r)
		code := errors.ValidationCode(conversionErr.Err.Error())
		former.AddFieldMessage(conversionErr.Key, i.handler.Msg(code))
		response.ExecuteTemplate(template, former)

	} else if multiErr, ok := err.(schema.MultiError); ok {
		template, former := i.handler.Response(r)
		internalError := false

		for _, err := range multiErr {
			if conversionErr, ok := err.(schema.ConversionError); ok {
				code := errors.ValidationCode(conversionErr.Err.Error())
				former.AddFieldMessage(conversionErr.Key, i.handler.Msg(code))
			} else {
				internalError = true
				i.handler.Logger().Error(errors.New(err))
			}
		}

		if internalError {
			response.ExecuteTemplate("500.html", data.NewInternalServerError(i.handler.HTTPId()))
		} else {
			response.ExecuteTemplate(template, former)
		}

	} else {
		i.handler.Logger().Error(errors.New(err))
		response.ExecuteTemplate("500.html", data.NewInternalServerError(i.handler.HTTPId()))
	}
}

func (i *POST) parse() {
	st := reflect.ValueOf(i.handler).Elem()

	for j := 0; j < st.NumField(); j++ {
		field := st.Type().Field(j)

		value := field.Tag.Get("request")
		if value == "all" || strings.Contains(value, "post") {
			i.handler.SetRequestValue(st.Field(j))
			break
		}
	}
}
