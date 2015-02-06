package interceptor

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/schema"
)

type requester interface {
	RequestValue() reflect.Value
	SetRequestValue(reflect.Value)
	Logger() *log.Logger
}

type Poster struct {
	trama.NopAJAXInterceptor
	handler requestResponser
}

func NewPoster(h requester) *Poster {
	return &Poster{handler: h}
}

func (i *Poster) Before(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	i.parse()

	if request := i.handler.RequestValue(); request.IsValid() {
		decoder := schema.NewDecoder()

		if err := r.ParseForm(); err != nil {
			i.handler.Logger().Error(err)
			// TODO
			return
		}

		if err := decoder.Decode(request.Addr().Interface(), r.Form); err != nil {
			i.handler.Logger().Error(err)
			// TODO
			return
		}
	}
}

func (i *Poster) parse() {
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
