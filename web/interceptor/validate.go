package interceptor

import (
	"net/http"
	"reflect"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/protocol"
)

type normalizer interface {
	Normalize()
}

type validator interface {
	Validate(bool) protocol.Translator
}

type messageRequester interface {
	RequestValue() reflect.Value
	SetMessage(protocol.Translator)
}

type Validate struct {
	trama.NopAJAXInterceptor
	handler messageRequester
}

func NewValidate(h messageRequester) *Validate {
	return &Validate{handler: h}
}

func (i *Validate) Before(w http.ResponseWriter, r *http.Request) {
	if !i.handler.RequestValue().IsValid() || i.handler.RequestValue().IsNil() {
		return
	}

	if normalize, ok := i.handler.RequestValue().Interface().(normalizer); ok {
		normalize.Normalize()
	}

	if validate, ok := i.handler.RequestValue().Interface().(validator); ok {
		checkMandatoryFields := false

		switch r.Method {
		case "PUT", "POST":
			checkMandatoryFields = true
		}

		msg := validate.Validate(checkMandatoryFields)
		if msg != nil {
			w.WriteHeader(http.StatusBadRequest)
			i.handler.SetMessage(msg)
		}
	}
}
