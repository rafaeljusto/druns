package interceptor

import (
	"br/core"
	"net/http"
	"reflect"
	"strings"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/web/templates/data"
	"github.com/rafaeljusto/schema"
)

type requester interface {
	RequestValue() reflect.Value
	SetRequestValue(reflect.Value)
	Logger() *log.Logger
	HTTPId() string
}

type POST struct {
	trama.NopWebInterceptor
	handler requester
}

func NewPOST(h requester) *POST {
	return &POST{handler: h}
}

func (i *POST) Before(response trama.Response, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	i.parse()

	if request := i.handler.RequestValue(); request.IsValid() {
		decoder := schema.NewDecoder()

		if err := r.ParseForm(); err != nil {
			i.handler.Logger().Error(core.NewError(err))
			response.ExecuteTemplate("500.html", data.NewInternalServerError(i.handler.HTTPId()))
			return
		}

		if request.CanAddr() {
			request = request.Addr()
		}

		if err := decoder.Decode(request.Interface(), r.Form); err != nil {
			i.handler.Logger().Error(core.NewError(err))
			response.ExecuteTemplate("500.html", data.NewInternalServerError(i.handler.HTTPId()))

			// TODO: The decode could found a validation problem in the certificate, so we should return
			// the form with the specific error instead of a "Internal Server Error"

			return
		}
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
